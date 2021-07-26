package matrix

import (
	"fmt"
	"strings"

	"git.mrcyjanek.net/mrcyjanek/selfbot/db"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var Client *mautrix.Client

// Login interactively and return an error or authtoken
func Login() error {
	if string(db.Get("meta.accesstoken")) == "" {
		switch Ask("Would you like to use accesstoken (2) or username/password combination (1)? Please input a number") {
		case "1":
			homeserver := Ask("Homeserver (eg. mrcyjanek.net)")
			username := Ask("Username (eg. cyjan)")
			password := Ask("Password (eg. ***** *** * **)")
			var err error
			Client, err = mautrix.NewClient(homeserver, "", "")
			if err != nil {
				return err
			}
			resp, err := Client.Login(&mautrix.ReqLogin{
				Type: "m.login.password",
				Identifier: mautrix.UserIdentifier{
					Type: mautrix.IdentifierTypeUser,
					User: username,
				},
				Password:         password,
				StoreCredentials: true,
			})
			if err != nil {
				return err
			}
			fmt.Println("accesstoken:", resp.AccessToken)
			fmt.Println("deviceid:", resp.DeviceID)
			fmt.Println("userid:", resp.UserID)
			db.Set("meta.homeserver", []byte(homeserver))        // mrcyjanek.net
			db.Set("meta.accesstoken", []byte(resp.AccessToken)) // ZaI......................JRY
			db.Set("meta.deviceid", []byte(resp.DeviceID))       // bt7s33Z2
			db.Set("meta.userid", []byte(resp.UserID))           // @cyjan:mrcyjanek.net
		case "2":
			homeserver := Ask("Homeserver (eg. mrcyjanek.net)")
			accesstoken := Ask("Access token (eg. dijbejbfaicpsbigwkcgauie)")
			username := Ask("Username (eg. cyjan)")
			var err error
			Client, err = mautrix.NewClient(homeserver, id.NewUserID(username, homeserver), accesstoken)
			if err != nil {
				return err
			}
		}
	}
	homeserver := string(db.Get("meta.homeserver"))
	accesstoken := string(db.Get("meta.accesstoken"))
	username := strings.Split(strings.Split(string(db.Get("meta.userid")), ":")[0], "@")[1]

	var err error
	Client, err = mautrix.NewClient(homeserver, id.NewUserID(username, homeserver), accesstoken)
	return err
}

func Ask(s string) (r string) {
	fmt.Println(s)
	fmt.Print(" > ")
	fmt.Scanf("%s", &r)
	return r
}