package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

type SepiaResponse struct {
	Total int `json:"total"`
	Data  []struct {
		ID                    int       `json:"id"`
		UUID                  string    `json:"uuid"`
		Score                 float32   `json:"score"`
		CreatedAt             time.Time `json:"createdAt"`
		UpdatedAt             time.Time `json:"updatedAt"`
		PublishedAt           time.Time `json:"publishedAt"`
		OriginallyPublishedAt string    `json:"originallyPublishedAt"`
		Category              struct {
			ID    int    `json:"id"`
			Label string `json:"label"`
		} `json:"category"`
		License struct {
			ID    int    `json:"id"`
			Label string `json:"unknown"`
		} `json:"license"`
		Language struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"language"`
		Privacy struct {
			ID    int    `json:"id"`
			Label string `json:"label"`
		} `json:"privacy"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Duration      int      `json:"duration"`
		Tags          []string `json:"tags"`
		ThumbnailPath string   `json:"thumbnailPath"`
		ThumbnailUrl  string   `json:"thumbnailUrl"`
		PreviewPath   string   `json:"previewPath"`
		PreviewUrl    string   `json:"previewUrl"`
		EmbedPath     string   `json:"embedPath"`
		EmbedUrl      string   `json:"embedUrl"`
		URL           string   `json:"url"`
		Views         int      `json:"views"`
		Likes         int      `json:"likes"`
		Dislikes      int      `json:"dislikes"`
		IsLive        bool     `json:"isLive"`
		NSFW          bool     `json:"nsfw"`
		Account       struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			URL         string `json:"url"`
			Host        string `json:"host"`
			Avatar      struct {
				URL       string    `json:"url"`
				Path      string    `json:"path"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"avatar"`
		} `json:"account"`
		Channel struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			URL         string `json:"url"`
			Host        string `json:"host"`
			Avatar      struct {
				URL       string    `json:"url"`
				Path      string    `json:"path"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"avatar"`
		} `json:'channel'`
	} `json:"data"`
}

var Event = event.EventMessage
var About = "!peertube 'mode (one/many)' 'results number (integer)' 'search query' - Return search results from sepia search"

func Handle(source mautrix.EventSource, evt *event.Event) {
	if evt.Sender != matrix.Client.UserID {
		return
	}
	args, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		//matrix.Client.SendText(evt.RoomID, err.Error())
		return
	}
	if len(args) >= 1 && args[0] == "!peertube" {
		matrix.Client.SendReaction(evt.RoomID, evt.ID, "processing...")
		if len(args) != 4 {
			matrix.Client.SendText(evt.RoomID, About)
			return
		}
		count, err := strconv.Atoi(args[2])
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		url := "https://sepiasearch.org/api/v1/search/videos?search=" + url.QueryEscape(args[3]) + "&boostLanguages[]=en&nsfw=true&start=0&count=" + strconv.Itoa(count) + "&sort=-match"
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
		var resp SepiaResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		switch args[1] {
		case "one":
			msg := "Search results for query: <code>" + args[3] + "</code>\n"
			for i := range resp.Data {
				j := resp.Data[i]
				msg += strconv.Itoa(i+1) + ". <a href=\"" + j.URL + "\">" + j.Name + "</a>\n"
			}
			matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(msg, false, true))
		case "many":
			for i := range resp.Data {
				j := resp.Data[i]
				msg := `🌐 <a href="%[1]s">%[2]s</a>
🗒️<i>%[3]s</i>
🐷<a href="%[4]s">%[5]s</a>
`
				// 1. URL to video
				// 2. Name
				// 3. Description
				// 4. Uploader URL
				// 5. Uploader name
				resp, err := matrix.Client.UploadLink(j.PreviewUrl)
				if err != nil {
					matrix.Client.SendText(evt.RoomID, "Unable to send preview:"+err.Error())
				} else {
					matrix.Client.SendImage(evt.RoomID, "", resp.ContentURI)
				}
				matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, format.RenderMarkdown(fmt.Sprintf(msg, j.URL, j.Name, j.Description, j.Account.URL, j.Account.Name), false, true))
			}
		}
		matrix.Client.RedactEvent(evt.RoomID, evt.ID, mautrix.ReqRedact{
			Reason: "[selfbot] This event already got processed",
		})
	}
}
