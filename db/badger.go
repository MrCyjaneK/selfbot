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

// Helpers

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

func fromByte(ID []byte, i interface{}) {
	var b bytes.Buffer
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&i)
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

/*
// Crypto store

type Store struct {
}

func (s *Store) Flush() error {
	return db.Sync()
}

func (s *Store) PutAccount(account *crypto.OlmAccount) {
	Set("meta.store.account", toByte(account))
}

func (s *Store) GetAccount() (account *crypto.OlmAccount) {
	fromByte(Get("meta.store.account"), &account)
	return account
}

func (s *Store) AddSession(key id.SenderKey, olmsession *crypto.OlmSession) error {
	// We don't need to check for errors here.
	l, _ := s.GetSessions(key)
	l = append(l, olmsession)
	Set("meta.store.session."+string(toByte(key)), toByte(l))
	var o crypto.OlmSession
	fromByte(Get("meta.store.session_latest."+string(toByte(key))), &o)
	log.Fatal("TODO: Implement storing of last session.", o.ID().String())
	return nil
}

func (s *Store) HasSession(key id.SenderKey) bool {
	return string(Get("meta.store.session."+string(toByte(key)))) != ""
}

func (s *Store) GetSessions(key id.SenderKey) (l crypto.OlmSessionList, e error) {
	fromByte(Get("meta.store.session."+string(toByte(key))), &l)
	return
}

func (s *Store) GetLatestSession(key id.SenderKey) (o *crypto.OlmSession, err error) {
	b := Get("meta.store.session_latest." + string(toByte(key)))
	if string(b) == "" {
		return o, errors.New("meta.store.session_latest." + string(toByte(key)) + " not found.")
	}
	fromByte(b, &o)
	return
}

func (s *Store) UpdateSession(key id.SenderKey, o *crypto.OlmSession) (err error) {
	l, _ := s.GetSessions(key)
	for i := range l {
		if l[i].ID().String() == o.ID().String() {
			l[i] = o
			break
		}
	}
	Set("meta.store.session."+string(toByte(key)), toByte(l))
	return
}

/*
	// PutGroupSession inserts an inbound Megolm session into the store. If an earlier withhold event has been inserted
	// with PutWithheldGroupSession, this call should replace that. However, PutWithheldGroupSession must not replace
	// sessions inserted with this call.
	PutGroupSession(id.RoomID, id.SenderKey, id.SessionID, *InboundGroupSession) error
	// GetGroupSession gets an inbound Megolm session from the store. If the group session has been withheld
	// (i.e. a room key withheld event has been saved with PutWithheldGroupSession), this should return the
	// ErrGroupSessionWithheld error. The caller may use GetWithheldGroupSession to find more details.
	GetGroupSession(id.RoomID, id.SenderKey, id.SessionID) (*InboundGroupSession, error)
	// PutWithheldGroupSession tells the store that a specific Megolm session was withheld.
	PutWithheldGroupSession(event.RoomKeyWithheldEventContent) error
	// GetWithheldGroupSession gets the event content that was previously inserted with PutWithheldGroupSession.
	GetWithheldGroupSession(id.RoomID, id.SenderKey, id.SessionID) (*event.RoomKeyWithheldEventContent, error)

	// GetGroupSessionsForRoom gets all the inbound Megolm sessions for a specific room. This is used for creating key
	// export files. Unlike GetGroupSession, this should not return any errors about withheld keys.
	GetGroupSessionsForRoom(id.RoomID) ([]*InboundGroupSession, error)
	// GetGroupSessionsForRoom gets all the inbound Megolm sessions in the store. This is used for creating key export
	// files. Unlike GetGroupSession, this should not return any errors about withheld keys.
	GetAllGroupSessions() ([]*InboundGroupSession, error)

	// AddOutboundGroupSession inserts the given outbound Megolm session into the store.
	//
	// The store should index inserted sessions by the RoomID field to support getting and removing sessions.
	// There will only be one outbound session per room ID at a time.
	AddOutboundGroupSession(*OutboundGroupSession) error
	// UpdateOutboundGroupSession updates the given outbound Megolm session in the store.
	UpdateOutboundGroupSession(*OutboundGroupSession) error
	// GetOutboundGroupSession gets the stored outbound Megolm session for the given room ID from the store.
	GetOutboundGroupSession(id.RoomID) (*OutboundGroupSession, error)
	// RemoveOutboundGroupSession removes the stored outbound Megolm session for the given room ID.
	RemoveOutboundGroupSession(id.RoomID) error

	// ValidateMessageIndex validates that the given message details aren't from a replay attack.
	//
	// Implementations should store a map from (senderKey, sessionID, index) to (eventID, timestamp), then use that map
	// to check whether or not the message index is valid:
	//
	// * If the map key doesn't exist, the given values should be stored and this should return true.
	// * If the map key exists and the stored values match the given values, this should return true.
	// * If the map key exists, but the stored values do not match the given values, this should return false.
	ValidateMessageIndex(senderKey id.SenderKey, sessionID id.SessionID, eventID id.EventID, index uint, timestamp int64) bool

	// GetDevices returns a map from device ID to DeviceIdentity containing all devices of a given user.
	GetDevices(id.UserID) (map[id.DeviceID]*DeviceIdentity, error)
	// GetDevice returns a specific device of a given user.
	GetDevice(id.UserID, id.DeviceID) (*DeviceIdentity, error)
	// PutDevice stores a single device for a user, replacing it if it exists already.
	PutDevice(id.UserID, *DeviceIdentity) error
	// PutDevices overrides the stored device list for the given user with the given list.
	PutDevices(id.UserID, map[id.DeviceID]*DeviceIdentity) error
	// FilterTrackedUsers returns a filtered version of the given list that only includes user IDs whose device lists
	// have been stored with PutDevices. A user is considered tracked even if the PutDevices list was empty.
	FilterTrackedUsers([]id.UserID) []id.UserID

	// PutCrossSigningKey stores a cross-signing key of some user along with its usage.
	PutCrossSigningKey(id.UserID, id.CrossSigningUsage, id.Ed25519) error
	// @carlod:mozilla.org retrieves a user's stored cross-signing keys.
	GetCrossSigningKeys(id.UserID) (map[id.CrossSigningUsage]id.Ed25519, error)
	// PutSignature stores a signature of a cross-signing or device key along with the signer's user ID and key.
	PutSignature(id.UserID, id.Ed25519, id.UserID, id.Ed25519, string) error
	// GetSignaturesForKeyBy returns the signatures for a cross-signing or device key by the given signer.
	GetSignaturesForKeyBy(id.UserID, id.Ed25519, id.UserID) (map[id.Ed25519]string, error)
	// IsKeySignedBy returns whether a cross-signing or device key is signed by the given signer.
	IsKeySignedBy(id.UserID, id.Ed25519, id.UserID, id.Ed25519) (bool, error)
	// DropSignaturesByKey deletes the signatures made by the given user and key from the store. It returns the number of signatures deleted.
	DropSignaturesByKey(id.UserID, id.Ed25519) (int64, error)
*/
