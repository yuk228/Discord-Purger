package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yuk228/disgolf"
)

type SearchResponse struct {
	AnalyticsID              string                   `json:"analytics_id"`
	DoingDeepHistoricalIndex bool                     `json:"doing_deep_historical_index"`
	TotalResults             int                      `json:"total_results"`
	Channels                 []interface{}            `json:"channels"`
	Messages                 [][]*discordgo.Message  `json:"messages"`
}

func HandlePurge2(prefix string) func(ctx *disgolf.MessageCtx) {
	return func(ctx *disgolf.MessageCtx) {
		switch len(ctx.Arguments) {
		case 0:
			ctx.Reply(fmt.Sprintf("```%spurge2 [channel_id] \n%spurge2 1234567891234567891```", prefix, prefix), false)

		case 1:
			channelID := ctx.Arguments[0]
			if len([]rune(channelID)) != 19 {
				ctx.Reply(fmt.Sprintln("length of channel id must be 19"), false)
				return
			}
			searchAndDelete(ctx, channelID)
		}
	}
}

func searchAndDelete(ctx *disgolf.MessageCtx, channelID string) {
	offset := 0
	totalDeleted := 0

	for {
		messages, hasMore, err := searchMessages(ctx, channelID, offset)
		if err != nil {
			log.Printf("```error searching messages: %s```", err.Error())
			return
		}

		if len(messages) == 0 {
			break
		}

		ownerIDs := strings.Split(os.Getenv("OWNER_IDS"), ",")
		for _, msg := range messages {
			if slices.Contains(ownerIDs, msg.Author.ID) && !slices.Contains([]discordgo.MessageType{3, 4, 5}, msg.Type) {
				err := ctx.ChannelMessageDelete(channelID, msg.ID)
				if err != nil {
					log.Printf("error deleting message %s: %s", msg.ID, err.Error())
				} else {
					totalDeleted++
					log.Printf("deleted: %s", msg.ID)
				}
				time.Sleep(time.Microsecond * 100)
			}
		}

		if !hasMore {
			break
		}

		offset += 25
		time.Sleep(1450 * time.Millisecond)
	}

	log.Printf("```deleted %d messages```", totalDeleted)
}


func searchMessages(ctx *disgolf.MessageCtx, channelID string, offset int) ([]*discordgo.Message, bool, error) {
	session := ctx.Session

	apiURL := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages/search", channelID)
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, false, err
	}

	params := url.Values{}
	params.Add("author_id", ctx.Message.Author.ID)
	params.Add("sort_by", "timestamp")
	params.Add("sort_order", "desc")
	params.Add("offset", strconv.Itoa(offset))
	reqURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, false, err
	}

	req.Header.Set("Authorization", session.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, false, err
	}

	var allMessages []*discordgo.Message
	for _, messageGroup := range searchResp.Messages {
		allMessages = append(allMessages, messageGroup...)

	}
	hasMore := len(allMessages) == 25 && offset+25 < searchResp.TotalResults

	return allMessages, hasMore, nil
}