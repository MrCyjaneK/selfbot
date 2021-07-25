package main

import (
	"log"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = "!ping - pong!"

func Handle(source mautrix.EventSource, evt *event.Event) {
	if evt.Sender != matrix.Client.UserID {
		return
	}
	if evt.Content.AsMessage().Body == "!ping" {
		_, err := matrix.Client.SendText(evt.RoomID, "pong!")
		if err != nil {
			log.Println(err)
		}
	}
}
