package main

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

//go:embed haregly.json
var contentj []byte
var Event = event.EventMessage
var About = []string{"!haregly <num> optional - Generates Haregly's messages"}
var Command = "!haregly"

func Handle(source mautrix.EventSource, evt *event.Event) {
	ok, args := matrix.ProcessMsg(*evt, Command)
	if !ok {
		return
	}

	var payload []string
	json.Unmarshal(contentj, &payload)
	max_num := len(payload)

	if len(args) == 1 {
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(max_num)
		var haregly_msg = "via @haregly: " + payload[number]
		matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
			MsgType: event.MsgText,
			Body:    " * " + haregly_msg,
			NewContent: &event.MessageEventContent{
				MsgType: event.MsgText,
				Body:    haregly_msg,
			},
			RelatesTo: &event.RelatesTo{
				Type:    event.RelReplace,
				EventID: evt.ID,
			},
		})
	}
	if len(args) == 2 {
		number, err := strconv.Atoi(args[1])
		if err != nil {
			matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
				MsgType: event.MsgText,
				Body:    " * via @haregly: error 414",
				NewContent: &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    "via @haregly: error 414",
				},
				RelatesTo: &event.RelatesTo{
					Type:    event.RelReplace,
					EventID: evt.ID,
				},
			})
		} else {
			if number <= max_num {
				haregly_msg := "via @haregly: " + payload[number]
				matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    " * " + haregly_msg,
					NewContent: &event.MessageEventContent{
						MsgType: event.MsgText,
						Body:    haregly_msg,
					},
					RelatesTo: &event.RelatesTo{
						Type:    event.RelReplace,
						EventID: evt.ID,
					},
				})
			} else {
				matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    " * via @haregly: error 414",
					NewContent: &event.MessageEventContent{
						MsgType: event.MsgText,
						Body:    "via @haregly: error 414",
					},
					RelatesTo: &event.RelatesTo{
						Type:    event.RelReplace,
						EventID: evt.ID,
					},
				})
			}
		}
		return
	}
}
