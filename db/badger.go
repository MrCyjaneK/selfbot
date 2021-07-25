package db

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"path"

	"github.com/dgraph-io/badger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var db *badger.DB
var DataDir = ""

func DBInit() {
	if DataDir == "" {
		log.Fatal("Please set DataDir!")
	}
	os.MkdirAll(path.Join(DataDir, "/db"), 0750)
	var err error
	db, err = badger.Open(badger.DefaultOptions(path.Join(DataDir, "/db")))
	if err != nil {
		log.Fatal(err)
	}
}

// Get value from db
func Get(key string) []byte {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if err.Error() != "Key not found\n" {
				log.Println(key, "not found")
				return nil
			}
			return err
		}
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return value
}

// Set value
func Set(key string, value []byte) {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Storer
type Storer struct {
	/*
		type Storer interface {
			SaveFilterID(userID id.UserID, filterID string)
			LoadFilterID(userID id.UserID) string
			SaveNextBatch(userID id.UserID, nextBatchToken string)
			LoadNextBatch(userID id.UserID) string
			SaveRoom(room *Room)
			LoadRoom(roomID id.RoomID) *Room
		}
	*/
}

func (s *Storer) SaveFilterID(userID id.UserID, filterID string) {
	Set("meta.storer.filters."+string(toByte(userID)), []byte(filterID))
}

func (s *Storer) LoadFilterID(userID id.UserID) string {
	return string(Get("meta.storer.filters." + string(toByte(userID))))
}

func (s *Storer) SaveNextBatch(userID id.UserID, nextBatchToken string) {
	Set("meta.storer.nextbatchtoken."+string(toByte(userID)), []byte(nextBatchToken))
}

func (s *Storer) LoadNextBatch(userID id.UserID) string {
	return string(Get("meta.storer.nextbatchtoken." + string(toByte(userID))))
}

func (s *Storer) SaveRoom(room *mautrix.Room) {
	Set("meta.storer.room."+string(toByte(room.ID)), toByte(&room))
}

func (s *Storer) LoadRoom(roomID id.RoomID) *mautrix.Room {
	var room mautrix.Room
	fromByte(Get("meta.storer.room."+string(toByte(room.ID))), &room)
	return &room
}

func toByte(i interface{}) (ID []byte) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(i)
	if err != nil {
		log.Fatal(err)
	}
	ID = b.Bytes()
	return
}

func fromByte(ID []byte, i interface{}) (j interface{}) {
	var b bytes.Buffer
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&i)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
