// +build ignore

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"os"
	"reflect"
	"text/template"
	"time"
)

type EventTypes = map[string]string

func main() {
	eventTypes := make(EventTypes, len(slackevents.EventsAPIInnerEventMapping))
	for name, event := range slackevents.EventsAPIInnerEventMapping {
		reflection := reflect.ValueOf(event)
		eventName := reflection.Type().Name()
		eventTypes[name] = eventName
	}

	f, err := os.Create("events_gen.go")
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to generate event callbacks")
	}
	defer f.Close()

	tpl, err := template.New("").Parse(eventsTemplate)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to parse events template")
	}
	err = tpl.Execute(f, struct {
		Timestamp  time.Time
		EventTypes EventTypes
	}{
		Timestamp:  time.Now(),
		EventTypes: eventTypes,
	})
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to execute events template")
	}
}

var eventsTemplate = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package slackbot

import "github.com/slack-go/slack/slackevents"
{{ range $key, $event := .EventTypes }}
type {{ $event }}Container struct {
	APIEvent slackevents.EventsAPIEvent
	Event slackevents.{{ $event }}
}

type {{ $event }}Callback = func(bot *Bot, c {{ $event }}Container)

func (b *Bot) Register{{ $event }}(callback {{ $event }}Callback) {
	b.RegisterEvent("{{ $key }}", func(bot *Bot, event slackevents.EventsAPIEvent) {
		e := event.InnerEvent.Data.(*slackevents.{{ $event }})
		callback(b, {{ $event }}Container{APIEvent: event, Event: *e})
	})
}
{{ end }}
`
