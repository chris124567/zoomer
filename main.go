package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/chris124567/zoomer/zoom"
)

func main() {
	meetingNumber := flag.String("meetingNumber", "", "Meeting number")
	meetingPassword := flag.String("password", "", "Meeting password")
	flag.Parse()

	// get keys from environment
	apiKey := os.Getenv("ZOOM_JWT_API_KEY")
	apiSecret := os.Getenv("ZOOM_JWT_API_SECRET")

	// create new session
	// meetingNumber, meetingPassword, username, hardware uuid (can be random but should be relatively constant or it will appear to zoom that you have many many many devices), proxy url, jwt api key, jwt api secret)
	session, err := zoom.NewZoomSession(*meetingNumber, *meetingPassword, "Bot", "ad8ffee7-d47c-4357-9ac8-965ed64e96fc", "", apiKey, apiSecret)
	if err != nil {
		panic(err)
	}
	// get the rwc token and other info needed to construct the websocket url for the meeting
	meetingInfo, cookieString, err := session.GetMeetingInfoData()
	if err != nil {
		panic(err)
	}

	// get the url for the websocket connection.  always pass false for the second parameter (its used internally to keep track of some parameters used for getting out of waiting rooms)
	websocketUrl, err := session.GetWebsocketUrl(meetingInfo, false)
	if err != nil {
		panic(err)
	}

	// the third argument is the "onmessage" function.  it will be triggered everytime the websocket client receives a message
	panic(session.MakeWebsocketConnection(websocketUrl, cookieString, func(session *zoom.ZoomSession, message zoom.Message) error {
		switch m := message.(type) {
		case *zoom.ConferenceRosterIndication:
			// if we get an indication that someone joined the meeting, welcome them
			for _, person := range m.Add {
				// don't welcome ourselves
				if person.ID != session.JoinInfo.UserID {
					// you could switch out EVERYONE_CHAT_ID with person.ID to private message them instead of sending the welcome to everyone
					session.SendChatMessage(zoom.EVERYONE_CHAT_ID, "Welcome to the meeting, "+string(person.Dn2)+"!")
				}
			}
			return nil
		case *zoom.ConferenceChatIndication:
			// respond to chats
			return handleChatMessage(session, m, string(m.Text))
		default:
			return nil
		}
	}))
}

func handleChatMessage(session *zoom.ZoomSession, body *zoom.ConferenceChatIndication, messageText string) error {
	// only respond to messages with this prefix
	const (
		messagePrefix = "++"
	)

	// takes commands of the form "++command argument1 argument2 ..."
	if !strings.HasPrefix(messageText, messagePrefix) {
		// this message is not for the bot
		return nil
	}

	words := strings.Fields(strings.TrimPrefix(messageText, messagePrefix))
	if len(words) < 1 {
		return errors.New("No command provided after prefix")
	}
	args := words[1:]

	switch words[0] {
	case "rename":
		if len(args) > 0 {
			session.RenameMe(strings.Join(args, " "))
		}
	case "mute":
		// if we get no arguments or "on", turn mute on
		if len(args) == 0 || args[0] == "on" {
			session.SetAudioMuted(true)
			session.SetVideoMuted(true)
		} else if args[0] == "off" {
			session.SetAudioMuted(false)
			session.SetVideoMuted(false)
		}
	case "screenshare":
		// if we get no arguments or "on", turn screenshare on
		if len(args) == 0 || args[0] == "on" {
			session.SetScreenShareMuted(false)
		} else if args[0] == "off" {
			session.SetScreenShareMuted(true)
		}
	case "chatlevel":
		// take the first argument, convert to integer and try to use that to set the room chat level
		if len(args) > 0 {
			if chatLevelInt, err := strconv.Atoi(args[0]); err == nil {
				session.SetChatLevel(chatLevelInt)
			}
		}
	default:
		// just echo the message it if its not code for anything
		session.SendChatMessage(body.DestNodeID, "I don't understand this message so I am echoing it: "+string(body.Text))
	}

	return nil
}
