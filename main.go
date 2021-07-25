package main

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"plugin"

	"git.mrcyjanek.net/mrcyjanek/selfbot/db"
	"git.mrcyjanek.net/mrcyjanek/selfbot/matrix"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var DataDir string

func main() {
	load()
	db.DBInit()
	err := matrix.Login()
	if err != nil {
		log.Fatal("Unable to login to matrix", err)
	}
	matrix.Client.Store = &db.Storer{}
	syncer := matrix.Client.Syncer.(*mautrix.DefaultSyncer)
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
			ev := pEvent.(*event.Type)
			syncer.OnEventType(*ev, pHandle.(func(mautrix.EventSource, *event.Event)))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
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
