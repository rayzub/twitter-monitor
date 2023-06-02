package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/webhook"
	"github.com/rayzub/twitter-monitor/src/twitter"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Context 			context.Context
	Monitor 			*twitter.Handler
	WHookClient 		webhook.Client
	BotClient  		 	bot.Client
	PingChan    		chan twitter.MonitorPing

	RequestChannelId 	string
}


func New(ctx context.Context) *Handler {
	pingChan := make(chan twitter.MonitorPing)
	monitor := twitter.New(pingChan)

	wClient, err := webhook.NewWithURL(os.Getenv("WEBHOOK"))

	if err != nil {
		return nil
	}

	bClient, err := disgo.New(os.Getenv("BOT_TOKEN"))

	if err != nil {
		return nil
	}

	handler := &Handler{
		Context:     		ctx,
		Monitor:     		monitor,
		WHookClient: 		wClient,
		BotClient:   		bClient,
		PingChan:    		pingChan,
		RequestChannelId:   os.Getenv("REQUEST_CHANNEL_ID"),
	}


	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ping := <-handler.PingChan
				if err := handler.SendWebhook(ping); err != nil {
					continue
				}
			}
		}
	}()

	
	return handler
}
func (h *Handler) SendWebhook(ping twitter.MonitorPing) error {
	return nil
}

func (h *Handler) HandleCommands(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID || message.ChannelID != h.RequestChannelId {
		return
	}

	if !strings.HasPrefix(message.Content, ".") {
		return
	}

	messageParts := strings.Split(message.Content, " ")
	command := strings.ToLower(messageParts[0])

	switch command {
	case ".add":
		if len(messageParts) > 1 {
			return
		}

		twitterHandle := strings.ToLower(messageParts[1])
		if err := h.addTwitterAccount(twitterHandle); err != nil {
			return
		}

		return
	case ".remove":
		if len(messageParts) > 1 {
			return
		}

		twitterHandle := strings.ToLower(messageParts[1])
		if err := h.removeTwitterAccount(twitterHandle); err != nil {
			return
		}

		return
	case ".list":	
		return			
	}
}


func (h *Handler) addTwitterAccount(handle string) error {
	accountId := h.Monitor.FetchTwitterID(handle)

	if accountId == 0 {
		return fmt.Errorf("error fetching twitter id for handle: %s", handle)
	}

	ctx, cancel := context.WithCancel(h.Context)
	syncChan := make(chan int)
	go twitter.MonitorTweets(h.Monitor, ctx, twitter.MonitorFilter{
		PositiveKeywords: []string{},
		NegativeKeywords: []string{},
		TwitterId: accountId,
		LatestTweetTS: time.Now().Unix(),
	})
	<-syncChan


	h.Monitor.CurrentMonitored = append(h.Monitor.CurrentMonitored, accountId)
	h.Monitor.MonitorKillMap[accountId] = cancel
	return nil
}

func (h *Handler) removeTwitterAccount(handle string) error {
	accountId := h.Monitor.FetchTwitterID(handle)

	if accountId == 0 {
		return fmt.Errorf("error fetching twitter id for handle: %s", handle)
	}

	cancel, ok := h.Monitor.MonitorKillMap[accountId]

	if !ok {
		return fmt.Errorf("%s is not being currently monitored", handle)
	}

	cancel()
	delete(h.Monitor.MonitorKillMap, accountId)
	accountIdIndx := slices.Index(h.Monitor.CurrentMonitored, accountId)
	h.Monitor.CurrentMonitored = slices.Delete(h.Monitor.CurrentMonitored, accountIdIndx, accountIdIndx+1)
	return nil
}

func (h *Handler) listMonitored() error {
	// return embed of accounts monitored
	return nil
}