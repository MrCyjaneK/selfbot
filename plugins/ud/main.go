package main

import (
	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage

func Handle(source mautrix.EventSource, evt *event.Event) {
	if evt.Sender != matrix.Client.UserID {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		return
	}
	if msgs[0] == "!ud" {
		if len(msgs) < 2 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!ud \"wat\"")
			return
		}

	}
}
