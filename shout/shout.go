package shout

import (
	"fmt"
	"strings"
	"time"

	"github.com/jurgen-kluft/hass-go/state"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

// Instance is our instant-messenger instance (currently Slack)
type Instance struct {
	viper *viper.Viper
	slack *slack.Client
}

// New creates a new instance of Slack
func New() (*Instance, error) {
	s := &Instance{}
	s.viper = viper.New()

	// Viper command-line package
	s.viper.SetConfigName("shout")   // name of config file (without extension)
	s.viper.AddConfigPath("config/") // optionally look for config in the working directory
	err := s.viper.ReadInConfig()    // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		return nil, err
	}

	s.slack = slack.New(s.viper.GetString("slack.key"))
	return s, nil
}

// PostMessage posts a message to a channel
func (s *Instance) postMessage(channel string, username string, msg string, pretext string, prebody string) {
	params := slack.PostMessageParameters{}
	params.Username = username
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    prebody,
	}
	params.Attachments = []slack.Attachment{attachment}
	_, timestamp, err := s.slack.PostMessage(channel, msg, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channel, timestamp)
}

func (s *Instance) Process(states *state.Domain) time.Duration {
	state := states.Get("shout")
	if len(state.Strings) > 0 {
		for name, body := range state.Strings {
			if strings.HasPrefix(name, "msg:") {
				parts := strings.SplitAfter(name, ":")
				if len(parts) == 2 && parts[0] == "msg:" {
					name = parts[1]
					channel := state.GetStringState(name+"."+"channel", "general")
					username := state.GetStringState(name+"."+"username", "bot")
					pretext := state.GetStringState(name+"."+"pretext", "...")
					prebody := state.GetStringState(name+"."+"prebody", "...")
					s.postMessage(channel, username, body, pretext, prebody)

				}
			}
		}
		state.Clear()
	}

	return 1 * time.Second
}
