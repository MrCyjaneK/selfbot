package main

import (
	"log"
	"embed"
	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"encoding/json"
	"strconv"
	"math/rand"
	"time"
)

//go:embed haregly.json
var f  embed.FS
var Event = event.EventMessage
var About = []string{"!haregly <num> optional - Generates Haregly's messages"}

func Handle(source mautrix.EventSource, evt *event.Event) {
	if !matrix.IsSelf(*evt) || matrix.IsOld(*evt) {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		return
	}
	if len(msgs) >= 1 && msgs[0] == "!haregly" {
		contentj, _ := f.ReadFile("haregly.json")
		var payload []interface{}
		json.Unmarshal(contentj, &payload)
		max_num := len(payload)

		if len(msgs) == 1 {
			rand.Seed(time.Now().UnixNano())
			number := rand.Intn(max_num)
			var haregly_msg = "via @haregly: "+payload[number].(string)
			matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
				MsgType: event.MsgText,
				Body:    " * "+haregly_msg,
				NewContent: &event.MessageEventContent{
					MsgType: event.MsgText,
					Body: haregly_msg,
				},
				RelatesTo: &event.RelatesTo{
					Type: event.RelReplace,
					EventID: evt.ID,
				},
			})
			if err != nil {
				log.Println(err)
			}
		}
		if len(msgs) == 2 {
		    number, err := strconv.Atoi(msgs[1])
		    if err != nil {
			    matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					    MsgType: event.MsgText,
					    Body:    " * via @haregly: error 414",
					    NewContent: &event.MessageEventContent{
						    MsgType: event.MsgText,
						    Body: "via @haregly: error 414",
					    },
					    RelatesTo: &event.RelatesTo{
						    Type: event.RelReplace,
						    EventID: evt.ID,
					    },
				    })
		    } else {
			    if (number <= max_num) {
				    haregly_msg := "via @haregly: "+ payload[number].(string)
				    matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    " * "+haregly_msg,
					NewContent: &event.MessageEventContent{
						MsgType: event.MsgText,
						Body: haregly_msg,
					},
					RelatesTo: &event.RelatesTo{
						Type: event.RelReplace,
						EventID: evt.ID,
					},
				    })
			    } else {
				    matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					    MsgType: event.MsgText,
					    Body:    " * via @haregly: error 414",
					    NewContent: &event.MessageEventContent{
						    MsgType: event.MsgText,
						    Body: "via @haregly: error 414",
					    },
					    RelatesTo: &event.RelatesTo{
						    Type: event.RelReplace,
						    EventID: evt.ID,
					    },
				    })
			    }
		    }
		    return
	    }
    }
}
