package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

var Event = event.EventMessage
var About = []string{"!ud 'word' - Get definition from urban dictionary"}
var Command = "!ud"

type ApiDefine struct {
	List []struct {
		Definition   string   `json:"definition"`
		Permalink    string   `json:"permalink"`
		ThumbsUp     int      `json:"thumbs_up"`
		ThumbsDown   int      `json:"thumbs_down"`
		SoundURL     []string `json:"sound_urls"`
		Author       string   `json:"author"`
		Word         string   `json:"word"`
		DefinitionID int      `json:"defid"`
		CurrentVote  string   `json:"current_vote"`
		WrittenOn    string   `json:"written_on"`
		Example      string   `json:"example"`
	} `json:"list"`
}

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
		wat := args[1]
		resp, err := http.Get("https://api.urbandictionary.com/v0/define?term=" + url.QueryEscape(wat))
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		defer resp.Body.Close()
		bhtml, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		var r ApiDefine
		err = json.Unmarshal(bhtml, &r)
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		if len(args) == 2 {
			args = append(args, "0")
		}
		msgsedit := 0
		for i := range args[2:] {
			j := args[2:][i]
			k, err := strconv.Atoi(j)
			if err != nil {
				matrix.Client.SendText(evt.RoomID, "Unable to get record for ID '"+j+"': "+err.Error())
				continue
			}
			if k >= len(r.List) {
				matrix.Client.SendText(evt.RoomID, "Unable to get record for ID '"+j+"': it doesn't exist")
				continue
			}
			d := r.List[k]
			message := fmt.Sprintf(`Term: <b>%[1]s</b>
🗒️ %[2]s

🎤 %[3]s
👍%[4]d 👎%[5]d`, d.Word, d.Definition, d.Example, d.ThumbsUp, d.ThumbsDown)
			if msgsedit == 0 {
				matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					Body:          " * " + format.RenderMarkdown(message, false, true).Body,
					Format:        format.RenderMarkdown(message, false, true).Format,
					FormattedBody: " * " + format.RenderMarkdown(message, false, true).FormattedBody,
					NewContent: &event.MessageEventContent{
						MsgType:       event.MsgText,
						Body:          format.RenderMarkdown(message, false, true).Body,
						Format:        format.RenderMarkdown(message, false, true).Format,
						FormattedBody: format.RenderMarkdown(message, false, true).FormattedBody,
					},
					RelatesTo: &event.RelatesTo{
						Type:    event.RelReplace,
						EventID: evt.ID,
					},
				})
				msgsedit++
			} else {
				matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(message, false, true))
			}
		}
	}
}
