package main

import (
	"context"
	"log"

	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"github.com/chromedp/chromedp"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = []string{"!screenshot 'url' 'Options'... - Take screenshot of a website."}

func Handle(source mautrix.EventSource, evt *event.Event) {
	if !matrix.IsSelf(*evt) || matrix.IsOld(*evt) {
		return
	}
	msgs, err := gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		return
	}
	if len(msgs) >= 1 && msgs[0] == "!screenshot" {
		if len(msgs) < 2 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `"+About[0]+"`")
			return
		}
		abc, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		url := msgs[1]

		var imageBuf []byte
		log.Println("1")
		if err := chromedp.Run(abc, screenshotTasks(url, &imageBuf)); err != nil {
			log.Fatal(err.Error())
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		log.Println(len(imageBuf))
		r, err := matrix.Client.UploadBytes(imageBuf, "image/png")
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			log.Fatal(err.Error())
			return
		}
		_, err = matrix.Client.SendImage(evt.RoomID, "Screenshot of "+msgs[1], r.ContentURI)

		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			log.Fatal(err.Error())
			return
		}
		matrix.Client.RedactEvent(evt.RoomID, evt.ID, mautrix.ReqRedact{
			Reason: "[selfbot] This event already got processed",
		})
	}
}

func screenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.FullScreenshot(imageBuf, 100),
	}
}
