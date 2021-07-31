package main

import (
	"context"

	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"github.com/chromedp/chromedp"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var Event = event.EventMessage
var About = []string{"!screenshot 'url' 'Options'... - Take screenshot of a website."}
var Command = "!screenshot"

func Handle(source mautrix.EventSource, evt *event.Event) {
	ok, args := matrix.ProcessMsg(*evt, Command)
	if !ok {
		return
	}

	if len(args) >= 1 && args[0] == Command {
		if len(args) < 2 {
			matrix.Client.SendText(evt.RoomID, "Please use the correct syntax, for example `"+About[0]+"`")
			return
		}
		abc, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		url := args[1]

		var imageBuf []byte
		if err := chromedp.Run(abc, screenshotTasks(url, &imageBuf)); err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		r, err := matrix.Client.UploadBytes(imageBuf, "image/png")
		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
			return
		}
		_, err = matrix.Client.SendImage(evt.RoomID, "Screenshot of "+url, r.ContentURI)

		if err != nil {
			matrix.Client.SendText(evt.RoomID, err.Error())
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
