package zoom

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"time"
	// "github.com/google/uuid"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (session *ZoomSession) GetWebsocketUrl(meetingInfo *MeetingInfo, wasInWaitingRoom bool) (string, error) {
	pingRwcServer := getRwgPingServer(meetingInfo)
	rwgInfo, err := session.getRwgPingData(meetingInfo, pingRwcServer)
	if err != nil {
		return "", err
	}

	if len(meetingInfo.Result.EncryptedRWC) < 1 {
		return "", errors.New("No RWC hosts found")
	}

	// query string for websocket url
	values := url.Values{}

	values.Set("rwcAuth", rwgInfo.RwcAuth)
	values.Set("dn2", base64.StdEncoding.EncodeToString([]byte(meetingInfo.Result.UserName)))
	values.Set("auth", meetingInfo.Result.Auth)
	values.Set("sign", meetingInfo.Result.Sign)
	values.Set("browser", USER_AGENT_SHORTHAND)
	values.Set("trackAuth", meetingInfo.Result.TrackAuth)
	values.Set("mid", meetingInfo.Result.Mid)
	values.Set("tid", meetingInfo.Result.Tid)
	values.Set("lang", "en")
	values.Set("ts", strconv.FormatInt(meetingInfo.Result.Ts, 10))
	// values.Set("ZM-CID", uuid.New().String()) // random uuid
	values.Set("ZM-CID", session.HardwareID.String()) // this is a hardware id.  you shouldnt have it change a bunch of times per ip or you will look highly suspicious
	values.Set("_ZM_MTG_TRACK_ID", "")
	values.Set("jscv", "1.8.6")
	// values.Set("jscv", "1.8.5")
	values.Set("fromNginx", "false")
	values.Set("zak", "")
	values.Set("mpwd", meetingInfo.Result.Password)
	values.Set("as_type", "1")

	/*
	   if you set as_type 2 all the chat messages will be encrypted.  i didnt do this beacuse it is much easier to just set as_type=1 than figuring out the whole aes-gcm sha256 mess of IVs, tags, etc

	   if you do want to do that, you may find this of use (add the following to https://github.com/zoom/sample-app-web/blob/master/Local/js/meeting.js):

	   var enc = new TextDecoder("ascii");
	   const oldCryptoEncrypt = window.crypto.subtle.encrypt;
	   function encShim(alg, key, data) {
	       console.log("Encrypt arguments: ", arguments);
	       console.log("Encrypt array buffer as string: ", enc.decode(data));
	       const exportKey = crypto.subtle.exportKey("raw", key);

	       console.log("Exported key", exportKey);
	       return oldCryptoEncrypt.apply(window.crypto.subtle, arguments);
	   }
	   const oldCryptoDecrypt = window.crypto.subtle.decrypt;
	   function decShim(alg, key, data) {
	       console.log("Decrypt arguments: ", arguments);
	       const exportKey = crypto.subtle.exportKey("raw", key);

	       console.log("Exported key", exportKey);
	       return oldCryptoDecrypt.apply(window.crypto.subtle, arguments);
	   }

	   const oldCryptoSign = window.crypto.subtle.sign;
	   function signShim(alg, key, data) {
	       console.log("Sign arguments: ", arguments);
	       const exportKey = crypto.subtle.exportKey("raw", key);

	       console.log("Exported key", exportKey);
	       return oldCryptoSign.apply(window.crypto.subtle, arguments);
	   }

	   const oldCryptoImportKey = window.crypto.subtle.importKey;
	   function importKeyShim(format,keyData,algorithm,extractable,keyUsages) {
	       console.log("Import key arguments: ", arguments);

	       return oldCryptoImportKey.apply(window.crypto.subtle, arguments);
	   }

	   Object.defineProperty(window.crypto.subtle, "encrypt", {value: encShim})
	   Object.defineProperty(window.crypto.subtle, "decrypt", {value: decShim})
	   Object.defineProperty(window.crypto.subtle, "sign", {value: signShim})
	   Object.defineProperty(window.crypto.subtle, "importKey", {value: importKeyShim})


	   and check the console.log output (this hooks all the subtlecrypto functions and logs args).  I attempted this but gave up because just setting as_type=1 to allow for plaintext was much much easier.
	   also check zoom webclient.js ("easyStore" may be especially of interest)

	   keys are derived from some kind of sha256 hmac of the message sent in this message (Evt: 7938; Seq: 4): {"encryptKey":"2mLZj3gEh5RRoHCBQ5n2LFe8ur4HqcKkwF77zxU0ilM"}
	*/
	// unknown
	values.Set("tk", "")
	values.Set("cfs", "0")
	// "opt" is a parameter to specify a meeting within a meeting, for instance breakout rooms or the main meeting in a meeting with waiting room enabled
	if wasInWaitingRoom {
		values.Set("opt", session.meetingOpt)
		values.Set("zoomid", session.JoinInfo.ZoomID)
		values.Set("participantID", strconv.Itoa(session.JoinInfo.ParticipantID))
	}

	return (&url.URL{
		Scheme:   "wss",
		Host:     rwgInfo.Rwg,
		Path:     fmt.Sprintf("/wc/api/%s", meetingInfo.Result.MeetingNumber),
		RawQuery: values.Encode(),
	}).String(), nil
}

type onMessage func(session *ZoomSession, message Message) error

func (session *ZoomSession) MakeWebsocketConnection(websocketUrl string, cookieString string, onMessageFunction onMessage) error {
	websocketHeaders := http.Header{}
	websocketHeaders.Set("Accept-Language", "en-US,en;q=0.9")
	websocketHeaders.Set("Cache-Control", "no-cache")
	websocketHeaders.Set("Origin", "http://localhost:9999")
	websocketHeaders.Set("Pragma", "no-cache")
	websocketHeaders.Set("User-Agent", USER_AGENT)
	websocketHeaders.Set("Cookie", cookieString)

	dialer := websocket.Dialer{
		// TODO: REMOVE -- DEV ONLY FOR CHARLES PROXY
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if session.ProxyURL != nil {
		dialer.Proxy = http.ProxyURL(session.ProxyURL)
	}

	connection, _, err := dialer.Dial(websocketUrl, websocketHeaders)
	if err != nil {
		return err
	}
	session.websocketConnection = connection

	defer connection.Close()

	done := make(chan struct{})

	var message *GenericZoomMessage
	var wasInWaitingRoom bool = false
	go func() {
		defer close(done)
		for {
			// reset struct
			message = &GenericZoomMessage{}

			err := connection.ReadJSON(&message)
			if err != nil {
				log.Print("failed to read:", err)
				return
			}
			log.Printf("Received message (Evt: %s = %d; Seq: %d): %s", MessageNumberToName[message.Evt], message.Evt, message.Seq, string(message.Body))

			switch message.Evt {
			/*
				if we receive a WS_CONF_JOIN_RES message (this is sent along with a bunch of other things when the websocket connection is established) will also store some info from the join response that is necessary for sending chats into the session state
				important that this is done before any other handling functions
			*/
			case WS_CONF_JOIN_RES:
				bodyData := JoinConferenceResponse{}
				err := json.Unmarshal(message.Body, &bodyData)
				if err != nil {
					log.Print("Failed to unmarshal json: %+v", err)
					return
				}
				session.JoinInfo = bodyData
			/* figure out whether we are in the waiting room or not */
			case WS_CONF_HOLD_CHANGE_INDICATION:
				bodyData := ConferenceHoldChangeIndication{}
				err := json.Unmarshal(message.Body, &bodyData)
				if err != nil {
					log.Print("Failed to unmarshal json: %+v", err)
					return
				}
				if bodyData.BHold == true {
					wasInWaitingRoom = true
				}
			/* get the opt for the waiting room */
			case WS_CONF_OPTION_INDICATION:
				if wasInWaitingRoom {
					bodyData := ConferenceOptionIndication{}
					err := json.Unmarshal(message.Body, &bodyData)
					if err != nil {
						log.Print("Failed to unmarshal json: %+v", err)
						return
					}
					session.meetingOpt = bodyData.Opt
				}
			}

			// dont run the user defined functions in the waiting room
			if !wasInWaitingRoom {
				// convert generic json message to go type
				m, err := GetMessageBody(message)
				if err != nil {
					log.Printf("Decoding message failed: %+v", err)
					continue
				}
				err = onMessageFunction(session, m)
				if err != nil {
					log.Printf("User defined function failed: %+v", err)
				}
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// zoom sends pings (aside from regular websocket ones) approximately every minute of the form "{"evt":0,"seq":74}"
	minutelyJsonPingTicker := time.NewTicker(60 * time.Second)
	defer minutelyJsonPingTicker.Stop()

	for {
		if wasInWaitingRoom { // get out of the select loop if we are in the waiting room otherwise the function is gonna return and not give us an opportunity to connect to the breakout room
			break
		}
		select {
		case <-minutelyJsonPingTicker.C:
			if wasInWaitingRoom { // get out of the select loop if we are in the waiting room otherwise the function is gonna return and not give us an opportunity to connect to the breakout room
				break
			}
			session.SendMessage(connection, WS_CONN_KEEPALIVE, nil)
			break
		case <-done:
			if wasInWaitingRoom { // get out of the select loop if we are in the waiting room otherwise the function is gonna return and not give us an opportunity to connect to the breakout room
				break
			}
			return nil
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := connection.WriteMessage(websocket.CloseMessage, []byte(""))
			if err != nil {
				return err
			}
			<-done
			return nil
		}
	}

	if wasInWaitingRoom {
		meetingInfo, cookieString, err := session.GetMeetingInfoData()
		if err != nil {
			log.Fatal(err)
		}
		websocketUrl, err := session.GetWebsocketUrl(meetingInfo, true)
		if err != nil {
			log.Fatal(err)
		}
		return session.MakeWebsocketConnection(websocketUrl, cookieString, onMessageFunction)
	}

	return nil
}
