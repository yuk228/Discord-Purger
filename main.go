package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/yuk228/disgolf"
)

var (
	prefix            string
	token             string
	max_message_count int
)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	prefix = os.Getenv("PREFIX")
	token = os.Getenv("TOKEN")
	max_message_count, err = strconv.Atoi(os.Getenv("MAX_MESSAGE_COUNT"))
	if err != nil {
		log.Printf("Error converting type")
	}

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

	bot.AddHandler(messageCreate)
	bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
		Prefixes:      []string{prefix},
		MentionPrefix: false,
	}))

	bot.State.MaxMessageCount = max_message_count

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("[Received message] %s from %s\n", m.Content, m.Author.Username)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	log.Printf("[Command detected] %s\n", m.Content)
}

func GetMessages(ctx *disgolf.MessageCtx, channelID string, limit int, m chan []*discordgo.Message) {
	messages, err := ctx.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		log.Printf("Error getting messages: %s", err)
	}
	m <- messages
}

func HandlePurge(ctx *disgolf.MessageCtx) {
	if len(ctx.Arguments) >= 1 {
		channelID := ctx.Arguments[0]
		// amount := ctx.Arguments[1]
		// limit := ctx.Arguments[2]

		parts := make(chan []*discordgo.Message)
		go GetMessages(ctx, channelID, max_message_count, parts)

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
				if err != nil {
					log.Fatal(err.Error())
				}
				time.Sleep(time.Microsecond * 100)
			}

		}

	} else {
		ctx.Reply(fmt.Sprintf("```%spurge [channel_id] [amount] [float(time)]\n%spurge 1234567891234567891 100 1.45```", prefix, prefix), false)
	}
}

func HasOwnerMiddleware(ctx *disgolf.MessageCtx) {
	if slices.Contains(strings.Split(os.Getenv("OWNER_IDS"), ","), ctx.Message.Author.ID) {
		ctx.Next()
	}
}
