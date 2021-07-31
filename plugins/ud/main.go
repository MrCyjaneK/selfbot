package main

import (
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = []string{"!ud 'word' - Get definition from urban dictionary"}
var Command = "!ud"

func Handle(source mautrix.EventSource, evt *event.Event) {
	ok, args := matrix.ProcessMsg(*evt, Command)
	if !ok {
		return
	}
	if len(args) >= 1 && args[0] == "!ud" {
		if len(args) < 2 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!ud \"wat\"")
			return
		}

	}
}
