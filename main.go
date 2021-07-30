package main

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"git.mrcyjanek.net/mrcyjanek/selfbot/db"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var DataDir string
var cmds []string

func main() {
	load()
	db.DBInit()
	err := matrix.Login()
	if err != nil {
		log.Fatal("Unable to login to matrix", err)
	}
	matrix.Client.Store = &db.Storer{}
	syncer := matrix.Client.Syncer.(*mautrix.DefaultSyncer)
	if string(db.Get("meta.loaded")) != "true" {
		log.Println("Doing an initial sync... The program will close after 30 seconds, please restart it afterwards.")
		go func() {
			err = matrix.Client.Sync()
			if err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Second * 30)
		db.Set("meta.loaded", []byte("true"))
		os.Exit(0)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			log.Println("Loading plugin:", info.Name())
			p, err := plugin.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			pHandle, err := p.Lookup("Handle")
			if err != nil {
				log.Fatal(err)
			}
			pEvent, err := p.Lookup("Event")
			if err != nil {
				log.Fatal(err)
			}
			pAbout, err := p.Lookup("About")
			if err != nil {
				log.Fatal(err)
			}
			abts := *pAbout.(*[]string)
			for i := range abts {
				cmds = append(cmds, abts[i])
			}
			ev := pEvent.(*event.Type)
			syncer.OnEventType(*ev, pHandle.(func(mautrix.EventSource, *event.Event)))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		if !matrix.IsSelf(*evt) || matrix.IsOld(*evt) {
			return
		}
		if evt.Content.AsMessage().Body == "!help" {
			matrix.Client.SendMessageEvent(evt.RoomID, event.EventMessage, &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    " * "+"List of available commands: \n  - "+strings.Join(cmds, "\n  - "),
					NewContent: &event.MessageEventContent{
						MsgType: event.MsgText,
						Body: "List of available commands: \n  - "+strings.Join(cmds, "\n  - "),
					},
					RelatesTo: &event.RelatesTo{
						Type: event.RelReplace,
						EventID: evt.ID,
					},
				    })
		}
	})
	err = matrix.Client.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func load() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	DataDir = path.Join(usr.HomeDir, ".SelfBot")
	db.DataDir = DataDir
}
