package im

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

// IM is our instant-messenger instance (currently Slack)
type IM struct {
	viper *viper.Viper
	slack *slack.Client
}

// New creates a new instance of Slack
func New() (*IM, error) {
	im := &IM{}
	im.viper = viper.New()

	// Viper command-line package
	im.viper.SetConfigName("hass-go-slack")        // name of config file (without extension)
	im.viper.AddConfigPath("$HOME/.hass-go-slack") // call multiple times to add many search paths
	im.viper.AddConfigPath(".")                    // optionally look for config in the working directory
	err := im.viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                                // Handle errors reading the config file
		return nil, err
	}

	im.slack = slack.New(im.viper.GetString("slack.key"))
	return im, nil
}

// PostMessage posts a message to a channel
func (im *IM) PostMessage(channel string, username string, msg string, pretext string, prebody string) {
	params := slack.PostMessageParameters{}
	params.Username = username
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    prebody,
	}
	params.Attachments = []slack.Attachment{attachment}
	_, timestamp, err := im.slack.PostMessage(channel, msg, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channel, timestamp)
}
