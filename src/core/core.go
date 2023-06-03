package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"github.com/gtuk/discordwebhook"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo/webhook"
	"github.com/rayzub/twitter-monitor/src/twitter"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Context 			context.Context
	Monitor 			*twitter.Handler
	WHookClient 		webhook.Client
	BotClient  		 	*discordgo.Session
	PingChan    		chan twitter.MonitorPing

	RequestChannelId 	string
}


func New(ctx context.Context) error {
	pingChan := make(chan twitter.MonitorPing)
	monitor := twitter.New(pingChan)
	wClient, err := webhook.NewWithURL(os.Getenv("WEBHOOK"))

	if err != nil {
		return fmt.Errorf("error creating webhook: %s", err.Error())
	}

	bClient, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))

	if err != nil {
		return fmt.Errorf("error creating bot client: %s", err.Error())
	}

	handler := &Handler{
		Context:     		ctx,
		Monitor:     		monitor,
		WHookClient: 		wClient,
		BotClient:   		bClient,
		PingChan:    		pingChan,
		RequestChannelId:   os.Getenv("REQUEST_CHANNEL_ID"),
	}

	handler.BotClient.AddHandler(handler.HandleCommands)
	handler.BotClient.Identify.Intents = discordgo.IntentsGuildMessages
	if err := handler.BotClient.Open(); err != nil {
		return fmt.Errorf("error opening bot websocket: %s", err.Error())
	}
	log.Printf("Bot %s currently online. Press CTRL-C to exit.", handler.BotClient.State.User.Username)

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


	return nil
}
func (h *Handler) SendWebhook(ping twitter.MonitorPing) error {
	username := "PRINT Twitter Monitor"
	twitterURL := fmt.Sprintf("https://twitter.com/%s", ping.Handle)
	fieldTitle := "Detected Tweet"
	message := discordwebhook.Message{
		Username: &username,
		Embeds: &[]discordwebhook.Embed{
			{
				Title:  nil,
				Url:   	&ping.URL,
				Author: &discordwebhook.Author{
					Name: &ping.Handle,
					IconUrl: &ping.Image,
					Url:  &twitterURL,
				},
				Fields: &[]discordwebhook.Field{
					{
						Name:  &fieldTitle,
						Value: &ping.Message,
					},
				},
			},
		},
	}
	if err := discordwebhook.SendMessage(os.Getenv("WEBHOOK"), message); err != nil {
		return err
	}
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
		if len(messageParts) < 2 {
			return
		}
		twitterHandle := strings.ToLower(messageParts[1])
		if err := h.addTwitterAccount(twitterHandle); err != nil {
			log.Println(err.Error())
			return
		}

		log.Printf("Added %s to monitor list.", twitterHandle)
		h.BotClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Added %s to monitor list.", twitterHandle))
		return
	case ".remove":
		if len(messageParts) < 2 {
			return
		}

		twitterHandle := strings.ToLower(messageParts[1])
		if err := h.removeTwitterAccount(twitterHandle); err != nil {
			return
		}

		log.Printf("Removed %s from monitor list.", twitterHandle)
		h.BotClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Removed %s from monitor list.", twitterHandle))
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
	go func(){
		for {
			select {
			case <-ctx.Done():
				return

			default:
				twitter.MonitorTweets(h.Monitor, twitter.MonitorFilter{
					PositiveKeywords: []string{},
					NegativeKeywords: []string{},
					TwitterId: accountId,
					LatestTweetTS: time.Now().Unix(),
				})
			}
		}
	}()
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