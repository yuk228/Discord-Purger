package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yuk228/disgolf"
)

var (
	prefix string
	token  string
)

func envLoad() {
	prefix = os.Getenv("PREFIX")
	token = os.Getenv("TOKEN")
	if prefix == "" || token == "" {
		log.Fatal("PREFIX or TOKEN is not set")
	}
}

func main() {
	envLoad()
	bot, err := disgolf.New(token)
	bot.Identify.Intents = discordgo.IntentsAll
	if err != nil {
		panic(err)
	}

	err = bot.Open()
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	fmt.Printf("[Logged in as %s]\n", bot.State.User.Username)
	fmt.Printf("Latency: %f ms\n\n", bot.HeartbeatLatency().Seconds())

	bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
		Prefixes:      []string{prefix},
		MentionPrefix: false,
	}))

	bot.State.MaxMessageCount = 100

	bot.Router.Register(&disgolf.Command{
		Name:               "purge",
		Description:        "purge self messages",
		MessageHandler:     disgolf.MessageHandlerFunc(HandlePurge),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware)},
	})

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot
}

func GetMessages(ctx *disgolf.MessageCtx, channelID string, limit int, m chan []*discordgo.Message) {
	var beforeID string = ""

	for {
		messages, err := ctx.ChannelMessages(channelID, limit, beforeID, "", "")
		if err != nil {
			ctx.Reply(fmt.Sprintf("```Error getting messages: %s```", err.Error()), false)
			break
		}
		if len(messages) == 0 {
			log.Println("No more messages")
			break
		}

		log.Printf("messages: %d, last messageID: %s", len(messages), messages[len(messages)-1].ID)

		m <- messages

		beforeID = messages[len(messages)-1].ID
		if len(messages) < limit {
			break
		}
	}
	close(m)
}

func HandlePurge(ctx *disgolf.MessageCtx) {
	if len(ctx.Arguments) >= 1 {
		channelID := ctx.Arguments[0]

		if len([]rune(channelID)) != 19 {
			ctx.Reply(fmt.Sprintln("length of channel id must be 19"), false)
			return
		}

		parts := make(chan []*discordgo.Message)
		go GetMessages(ctx, channelID, 100, parts)

		for messages := range parts {
			var target_msgs []string
			for _, msg := range messages {

				// []discordgo.MessageType{MessageTypeCall, MessageTypeChannelNameChange, MessageTypeChannelIconChange}
				// これらのメッセージは削除不可な為除外 (groupでのみ発生)
				if slices.Contains(strings.Split(os.Getenv("OWNER_IDS"), ","), msg.Author.ID) && !slices.Contains([]discordgo.MessageType{3, 4, 5}, msg.Type) {
					target_msgs = append(target_msgs, msg.ID)
				}
			}
			for _, target_msg := range target_msgs {
				err := ctx.ChannelMessageDelete(channelID, target_msg)
				log.Printf("%s", target_msg)
				if err != nil {
					log.Fatal(err.Error())
				}
				time.Sleep(time.Microsecond * 100)
			}
		}
	} else {
		ctx.Reply(fmt.Sprintf("```%spurge [channel_id]\n%spurge 1234567891234567891```", prefix, prefix), false)
	}
}

func HasOwnerMiddleware(ctx *disgolf.MessageCtx) {
	if slices.Contains(strings.Split(os.Getenv("OWNER_IDS"), ","), ctx.Message.Author.ID) {
		ctx.Next()
	}
}
