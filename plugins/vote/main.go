package main

import (
	"log"

	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

var Event = event.EventMessage
var About = "!vote 'Question' 'Options'... - Create a voting poll"

type WikiResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			NS      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func Handle(source mautrix.EventSource, evt *event.Event) {
	if evt.Sender != matrix.Client.UserID {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		return
	}
	if len(msgs) >= 1 && msgs[0] == "!vote" {
		if len(msgs) < 3 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!vote \"Question\" \"Options\"...")
			return
		}
		r, err := matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(msgs[1], false, true))
		if err != nil {
			log.Println(err)
		}
		for code := range msgs[2:] {
			thing := msgs[code+2]
			matrix.Client.SendReaction(evt.RoomID, r.EventID, thing)
		}
		matrix.Client.RedactEvent(evt.RoomID, evt.ID, mautrix.ReqRedact{
			Reason: "[selfbot] This event already got processed",
		})
	}
}
