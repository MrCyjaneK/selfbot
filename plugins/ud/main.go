package main

import (
	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = []string{"!ud 'word' - Get definition from urban dictionary"}

func Handle(source mautrix.EventSource, evt *event.Event) {
	if !matrix.IsSelf(*evt) || matrix.IsOld(*evt) {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		//matrix.Client.SendText(evt.RoomID, err.Error())
		return
	}
	if len(msgs) >= 1 && msgs[0] == "!ud" {
		if len(msgs) < 2 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!ud \"wat\"")
			return
		}

	}
}
