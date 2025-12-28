package commands

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yuk228/disgolf"
)

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

func HandlePurge(prefix string) func(ctx *disgolf.MessageCtx) {
	return func(ctx *disgolf.MessageCtx) {
		if len(ctx.Arguments) >= 1 {
			channelID := ctx.Arguments[0]

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
						log.Printf("%s", err.Error())
					}
					time.Sleep(time.Microsecond * 100)
				}
			}
		} else {
			ctx.Reply(fmt.Sprintf("```%spurge [channel_id]\n%spurge 1234567891234567891```", prefix, prefix), false)
		}
	}
}

