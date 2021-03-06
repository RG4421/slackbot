[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/bushelpowered/slackbot/CI/master)](https://github.com/bushelpowered/slackbot/actions?query=branch%3Amaster+workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/bushelpowered/slackbot)](https://goreportcard.com/report/github.com/bushelpowered/slackbot)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/bushelpowered/slackbot?label=stable)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/bushelpowered/slackbot)](https://pkg.go.dev/mod/github.com/bushelpowered/slackbot)

# Slackbot
A simple callback-based framework for quickly building Slack Bots in Go built on https://github.com/slack-go/slack.

It supports:
* [Slash Commands](https://api.slack.com/interactivity/slash-commands)
* [Events API](https://api.slack.com/events-api)
* [Block Kit Interactivity](https://api.slack.com/block-kit/interactivity)
* [Shortcuts](https://api.slack.com/interactivity/shortcuts)
* [Option Load URL](https://api.slack.com/legacy/message-menus#adding-menus-to-messages__populate-message-menus-dynamically__options-load-url)

## Install

```console
go get github.com/bushelpowered/slackbot
```

## Basic Example

```go
package main

import (
	"github.com/bushelpowered/slackbot"
	"github.com/slack-go/slack"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Boot a bot with a slash command that posts "Hello World!" and message listener.
func main() {
	bot := slackbot.NewBot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_SIGNING_SECRET"))

	// register a slash command
	bot.RegisterCommand("myslashcommand", func(bot *slackbot.Bot, command slack.SlashCommand) *slack.Msg {
		log.Println(command.Text)
		return &slack.Msg{Text: "Hello World!"} // return nil for no reply
	})

	// register a message event
	bot.RegisterMessageEvent(func(bot *slackbot.Bot, c slackbot.MessageEventContainer) {
		log.Println(c.Event.Text)
	})

	// boot the bot
	err := bot.Boot(":8000")
	if err != nil {
		log.Println(err)
		return
	}
	defer bot.Shutdown(time.Second * 10)

	// wait for exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down...")
}
```

See more examples in [examples](examples).