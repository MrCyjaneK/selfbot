package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

var Event = event.EventMessage
var About = []string{"!wiki 'langcode (en)' 'search (Stack Overflow)'... "}

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

var msgformat = `🌐<b>%s</b>
🗒️<i>%s</i>`

func Handle(source mautrix.EventSource, evt *event.Event) {
	if !matrix.IsSelf(*evt) || matrix.IsOld(*evt) {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		return
	}
	if len(msgs) >= 1 && msgs[0] == "!wiki" {
		matrix.Client.SendReaction(evt.RoomID, evt.ID, "processing...")
		if len(msgs) < 3 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `!wiki \"langcode (en)\" \"search (Stack Overflow)\"")
			return
		}
		args_count := 0
		for code := range msgs[2:] {
			url := "https://" + url.QueryEscape(msgs[1]) + ".wikipedia.org/w/api.php?action=query&format=json&prop=extracts&titles=" + url.QueryEscape(msgs[code+2]) + "&exsentences=5&exlimit=1&exintro=1&explaintext=1"
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				matrix.Client.SendText(evt.RoomID, err.Error())
				return
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				matrix.Client.SendText(evt.RoomID, err.Error())
				return
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				matrix.Client.SendText(evt.RoomID, err.Error())
				return
			}
			var wikiresp WikiResponse
			err = json.Unmarshal(body, &wikiresp)
			if err != nil {
				matrix.Client.SendText(evt.RoomID, err.Error())
				return
			}
			for i := range wikiresp.Query.Pages {
				j := wikiresp.Query.Pages[i]
				if j.Extract == "" {
					j.Extract = "We were unable to define this term."
				}
				message := fmt.Sprintf(msgformat, j.Title, j.Extract)
				if args_count == 0 {
					matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
						Body: " * "+format.RenderMarkdown(message, false, true).Body,
						Format: format.RenderMarkdown(message, false, true).Format,
						FormattedBody: " * "+format.RenderMarkdown(message, false, true).FormattedBody,
						NewContent: &event.MessageEventContent{
							MsgType: event.MsgText,
							Body: format.RenderMarkdown(message, false, true).Body,
							Format: format.RenderMarkdown(message, false, true).Format,
							FormattedBody: format.RenderMarkdown(message, false, true).FormattedBody,
						},
						RelatesTo: &event.RelatesTo{
							Type: event.RelReplace,
							EventID: evt.ID,
						},
					})
				} else {
					matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(message, false, true))
				}
			}
			args_count ++
		}
	}
}
