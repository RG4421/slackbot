package main

import (
	"github.com/bushelpowered/slackbot"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Boot a bot with a slash command that echos Hello World!
func main() {
	bot := slackbot.NewBot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_SIGNING_SECRET"))

	// register command
	bot.RegisterCommand("test", exampleCommandCallback)

	// boot the bot
	err := bot.Boot(":8000")
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to start bot")
		return
	}
	defer bot.Shutdown(time.Second * 10)

	// wait for exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Infoln("Shutting down...")
}

func exampleCommandCallback(bot *slackbot.Bot, command slack.SlashCommand) *slack.Msg {
	logrus.Info(command)
	return &slack.Msg{Text: "Hello World!"} // return nil for no reply
}
