package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const twitchUsersURL = "https://api.twitch.tv/helix/users"
const twitchStreamsURL = "https://api.twitch.tv/helix/streams"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchWatchUsername := os.Getenv("TWITCH_WATCH_USERNAME")
	twitchClient := &twitchClient{
		client:   &http.Client{},
		clientID: twitchClientID,
	}

	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")
	discordNotificationChannelID := os.Getenv("DISCORD_NOTIFICATION_CHANNEL_ID")
	discord, err := discordgo.New("Bot " + discordBotToken)

	isStreaming := false
	discord.AddHandler(ready(&isStreaming, twitchClient, twitchWatchUsername, discordNotificationChannelID))
	discord.AddHandler(messageCreate(&isStreaming))

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Rumi is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func ready(isStreaming *bool, twitchClient *twitchClient, twitchWatchUsername string, discordNotificationID string) func(*discordgo.Session, *discordgo.Ready) {
	user := twitchClient.getUserInfo(twitchWatchUsername)
	if user == nil {
		panic(fmt.Sprintf("Invalid twitch user %s\n", twitchWatchUsername))
	}
	return func(s *discordgo.Session, event *discordgo.Ready) {
		go func() {
			for {
				streamData := twitchClient.getStreams(user.Data[0].ID)
				if streamData != nil && len(streamData.Data) > 0 {
					// Proper way would probably be to use channels, but I'm being lazy
					if !*isStreaming {
						*isStreaming = true
						s.ChannelMessageSend(discordNotificationID, fmt.Sprintf("Hey @here! %s is live! Check her out at https://twitch.tv/%s", user.Data[0].DisplayName, user.Data[0].Login))
						time.Sleep(10 * time.Second)
					}
				} else if *isStreaming {
					*isStreaming = false
					s.ChannelMessageSend(discordNotificationID, "The stream has ended for now")
				}
				time.Sleep(2 * time.Second)
			}
		}()

		s.UpdateStatus(0, "!rumi")
	}
}

func messageCreate(isStreaming *bool) func (s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// check if the message is "!"
		if strings.HasPrefix(strings.TrimSpace(m.Content), "!rumi") {
			action := strings.TrimSpace(strings.TrimPrefix(m.Content, "!rumi"))
			if action == "ping" {
				s.ChannelMessageSend(m.ChannelID, "Pong!")
			} else if action == "live" {
				if *isStreaming {
					s.ChannelMessageSend(m.ChannelID, "Live now")
				} else {
					s.ChannelMessageSend(m.ChannelID, "Currently offline")
				}
			}
		}
	}
}
type twitchClient struct {
	client *http.Client
	clientID string
}

type streamData struct {

	ID string `json:"id"`
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	GameID string `json:"game_id"`
	Type string `json:"type"`
	Title string `json:"title"`
	ViewerCount int `json:"viewer_count"`
	StartedAt time.Time `json:"started_at"`
	Language string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type twitchStreamResponse struct {
	Data []streamData `json:"data"`
}

func (c *twitchClient) getTwitchData(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("error creating request %v\n", err)
		return nil, err
	}
	req.Header.Add("Client-ID", c.clientID)
	res, err := c.client.Do(req)
	if err != nil {
		fmt.Errorf("error getting streams %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("error reading body %v\n", err)
	}
	return body, nil
}

func (c *twitchClient) getStreams(userID string) *twitchStreamResponse {
	var data twitchStreamResponse
	body, err := c.getTwitchData(fmt.Sprintf("%s?user_id=%s", twitchStreamsURL, userID))
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}
	return &data
}

type userData struct {

	ID string `json:"id"`
	Login string `json:"login"`
	DisplayName string `json:"display_name"`
	Type string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount int `json:"view_count"`
}


type twitchUserResponse struct {
	Data []userData `json:"data"`
}

func (c *twitchClient) getUserInfo(username string) *twitchUserResponse {
	data := twitchUserResponse{}
	body, err := c.getTwitchData(fmt.Sprintf("%s?login=%s", twitchUsersURL, username))
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}
	return &data
}
//
//{"data":
//	[
//		{
//			"id":"223957213",
//			"login":"lilylefae",
//			"display_name":"lilylefae",
//			"type":"",
//			"broadcaster_type":"",
//			"description":"Hey! I'm Lily and I'm a brazilian streamer living in the USA!",
//			"profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/b88e5be0-4d72-4634-b908-8de811536a7b-profile_image-300x300.png",
//			"offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/a9d4fc08-0930-4098-9b4d-201aea98d212-channel_offline_image-1920x1080.png",
//			"view_count":7961
//		}
//	]
//}