package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yuk228/disgolf"

	"github.com/yuk228/Discord-Purger/commands"
	"github.com/yuk228/Discord-Purger/middleware"
)

var (
	prefix string
	token  string
)

func init() {
	prefix = os.Getenv("PREFIX")
	token = os.Getenv("TOKEN")
	if prefix == "" || token == "" {
		log.Fatal("PREFIX or TOKEN is not set")
	}
}

func main() {
	bot, err := disgolf.New(token)
	if err != nil {
		log.Fatalf("error creating bot session: %v", err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
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
		MessageHandler:     disgolf.MessageHandlerFunc(commands.HandlePurge(prefix)),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(middleware.HasOwnerMiddleware)},
	})

	bot.Router.Register(&disgolf.Command{
		Name:               "purge2",
		Description:        "purge self messages using search api",
		MessageHandler:     disgolf.MessageHandlerFunc(commands.HandlePurge2(prefix)),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(middleware.HasOwnerMiddleware)},
	})

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot
}
