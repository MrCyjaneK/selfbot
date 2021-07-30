package matrix

import (
	"time"

	"maunium.net/go/mautrix/event"
)

// IsOld return true if event is older than 5 seconds.
//TODO: make this configurable
func IsOld(evt event.Event) bool {
	return -5*1000 > (evt.Timestamp - time.Now().Unix()*1000)
}

func IsSelf(evt event.Event) bool {
	return evt.Sender == Client.UserID
}
