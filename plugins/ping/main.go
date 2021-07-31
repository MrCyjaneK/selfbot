package main

import (
	"log"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = []string{"!ping - pong!"}
var Command = "!ping"

func Handle(source mautrix.EventSource, evt *event.Event) {
	ok, args := matrix.ProcessMsg(*evt, Command)
	if !ok {
		return
	}

	if args[0] == Command {
		_, err := matrix.Client.SendText(evt.RoomID, "pong!")
		if err != nil {
			log.Println(err)
		}
	}
}
