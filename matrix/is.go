package matrix

import (
	"strings"
	"time"

	gosh "git.mrcyjanek.net/mrcyjanek/gosh/_core"

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

func ProcessMsg(evt event.Event, Command string) (ok bool, args []string) {
	if !IsSelf(evt) || IsOld(evt) {
		return false, args
	}
	var err error
	args, err = gosh.Split(evt.Content.AsMessage().Body)
	if err != nil {
		var msgs = strings.Split(evt.Content.AsMessage().Body, " ")
		if len(msgs) >= 1 && msgs[0] == Command {
			Client.SendText(evt.RoomID, err.Error())
		}
		return false, args
	}
	return len(args) >= 1 && args[0] == Command, args
}
