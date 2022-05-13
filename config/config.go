package config

import (
	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Debugging bool

// InitializeConfiguration inits the base configuration of commentron
func InitializeConfiguration() {
	if viper.GetBool("debugmode") {
		Debugging = true
		logrus.SetLevel(logrus.DebugLevel)
	}
	if viper.GetBool("tracemode") {
		Debugging = true
		logrus.SetLevel(logrus.TraceLevel)
	}
	if viper.GetString("slack_hook") != "" {
		initSlack()
	}
}

// initSlack initializes the slack connection and posts info level or greater to the set channel.
func initSlack() {
	slackURL := viper.GetString("slack_hook")
	slackChannel := viper.GetString("slack_channel")
	if slackURL != "" && slackChannel != "" {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        slackURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
			Channel:        slackChannel,
			IconEmoji:      ":prism:",
			Username:       "Mirage",
		})
	}
}
