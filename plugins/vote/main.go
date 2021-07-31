package main

import (
	"log"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

var Event = event.EventMessage
var About = []string{"!vote 'Question' 'Options'... - Create a voting poll"}
var Command = "!vote"

func Handle(source mautrix.EventSource, evt *event.Event) {
	ok, args := matrix.ProcessMsg(*evt, Command)
	if !ok {
		return
	}
	if len(args) >= 1 && args[0] == Command {
		if len(args) < 3 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!vote \"Question\" \"Options\"...")
			return
		}
		r, err := matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(args[1], false, true))
		if err != nil {
			log.Println(err)
		}
		for code := range args[2:] {
			thing := args[code+2]
			matrix.Client.SendReaction(evt.RoomID, r.EventID, thing)
		}
		matrix.Client.RedactEvent(evt.RoomID, evt.ID, mautrix.ReqRedact{
			Reason: "[selfbot] This event already got processed",
		})
	}
}
