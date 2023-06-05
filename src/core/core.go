package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gtuk/discordwebhook"
	"github.com/rayzub/twitter-monitor/src/twitter"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Context 			context.Context
	Monitor 			*twitter.Handler
	BotClient  		 	*discordgo.Session
	Logger 				*zap.Logger
	PingChan    		chan twitter.MonitorPing

	RequestChannelId 	string
}


func New(ctx context.Context, logger *zap.Logger) error {
	pingChan := make(chan twitter.MonitorPing)
	monitor := twitter.New(pingChan)
	bClient, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))

	if err != nil {
		return fmt.Errorf("error creating bot client: %s", err.Error())
	}

	handler := &Handler{
		Context:     		ctx,
		Monitor:     		monitor,
		Logger:             logger,
		BotClient:   		bClient,
		PingChan:    		pingChan,
		RequestChannelId:   os.Getenv("REQUEST_CHANNEL_ID"),
	}

	handler.BotClient.AddHandler(handler.HandleCommands)
	handler.BotClient.Identify.Intents = discordgo.IntentsGuildMessages
	if err := handler.BotClient.Open(); err != nil {
		return fmt.Errorf("error opening bot websocket: %s", err.Error())
	}
	handler.Logger.Info(fmt.Sprintf("Bot %s currently online. Press CTRL-C to exit.", handler.BotClient.State.User.Username))
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
	title := fmt.Sprintf("New Tweet - %s", ping.Handle)
	fieldTitle := "Detected Tweet"
	linkTitles := "Links"
	linkMessage := fmt.Sprintf("[Twitter Link](%s)", ping.URL)
	footerText := "Made by xyz#0004"
	inline := true
	message := discordwebhook.Message{
		Username: &username,
		Embeds: &[]discordwebhook.Embed{
			{
				Title:  &title,
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
					{
						Name: &linkTitles,
						Value: &linkMessage,
					},
				},
				Footer: &discordwebhook.Footer{
					Text: &footerText,
				},
			},
		},
	}

	if ping.MessageImage != "" {
		embed := *message.Embeds
		embed[0].Image = &discordwebhook.Image{
			Url: &ping.MessageImage,
		}
	}

	if len(ping.ParsedData) > 0 {
		embed := *message.Embeds
		for _, extraData := range ping.ParsedData {
			*embed[0].Fields = append(*embed[0].Fields, discordwebhook.Field{
				Name: &extraData.Title,
				Value: &extraData.Value,
				Inline: &inline,
			})
		}
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

		twitterHandles := strings.Split(strings.ToLower(messageParts[1]), ",")
		for _, twitterHandle := range twitterHandles {
			if strings.HasPrefix(twitterHandle, "https://") {
				splitTwitterURL := strings.Split(twitterHandle, "https://twitter.com/")
				twitterHandle = splitTwitterURL[1]
			}
	
			if err := h.addTwitterAccount(twitterHandle); err != nil {
				h.Logger.Error(err.Error())
				return
			}
	
			h.Logger.Info(fmt.Sprintf("Added %s to monitor list.", twitterHandle))
			h.BotClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Added %s to monitor list.", twitterHandle))
		}
		return
	case ".remove":
		if len(messageParts) < 2 {
			return
		}

		twitterHandles := strings.Split(strings.ToLower(messageParts[1]), ",")		
		for _, twitterHandle := range twitterHandles {
			if err := h.removeTwitterAccount(twitterHandle); err != nil {
				return
			}

			h.Logger.Info(fmt.Sprintf("Removed %s from monitor list.", twitterHandle))
			h.BotClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Removed %s from monitor list.", twitterHandle))
		}
		return
	case ".list":
		h.BotClient.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Title: "Current Monitored Twitters",
			Description: strings.Join(h.Monitor.CurrentMonitored, "\n"),
		})	
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
		twitterFilter := &twitter.MonitorFilter{
			PositiveKeywords: []string{},
			NegativeKeywords: []string{},
			TwitterId: accountId,
			LatestTweetTS: time.Now().Unix(),
		}

		for {
			select {
			case <-ctx.Done():
				return

			default:
				twitter.MonitorTweets(h.Monitor, twitterFilter)
			}
		}
	}()
	h.Monitor.CurrentMonitored = append(h.Monitor.CurrentMonitored, handle)
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
	accountIdIndx := slices.Index(h.Monitor.CurrentMonitored, handle)
	h.Monitor.CurrentMonitored = slices.Delete(h.Monitor.CurrentMonitored, accountIdIndx, accountIdIndx+1)
	return nil
}

