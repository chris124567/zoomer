package zoom

import (
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	PING             = 0x00
	RTP              = 0x4D
	RTCP             = 0x4E
	AES_GCM_IV_VALUE = 0x42
)

type ZoomSharingStream struct {
	mutex sync.Mutex
	recv  *websocket.Conn
	send  *websocket.Conn

	decoder *ZoomRtpDecoder

	mySecretNonce []byte
	sendQueue     chan []byte
}

func CreateZoomSharingStream(session *ZoomSession) (*ZoomSharingStream, error) {
	if session.JoinInfo == nil {
		return nil, errors.New("Zoom session does not have valid JoinInfo")
	}

	if session.RwgInfo == nil {
		return nil, errors.New("Zoom session does not have valid RwgInfo")
	}

	if session.JoinInfo.ZoomID == "" {
		return nil, errors.New("Zoom session does not have valid ZoomID")
	}

	secretNonce, err := ZoomEscapedBase64Decode(session.JoinInfo.ZoomID)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("type", "s")
	values.Set("cid", session.JoinInfo.ConID)

	baseUrl := &url.URL{
		Scheme:   "wss",
		Host:     session.RwgInfo.Rwg,
		Path:     "/wc/media/" + session.MeetingNumber,
		RawQuery: values.Encode(),
	}

	recv, err := createWebsocket("SharingRecv", baseUrl.String()+"&mode=1")
	if err != nil {
		return nil, err
	}

	send, err := createWebsocket("SharingSend", baseUrl.String()+"&mode=2")
	if err != nil {
		return nil, err
	}

	final := &ZoomSharingStream{
		recv:          recv,
		send:          send,
		decoder:       NewZoomRtpDecoder(session.ParticipantRoster),
		mySecretNonce: secretNonce,
	}

	go final.StartReceiveChannel()

	return final, nil
}

func Recorder() (io.WriteCloser, error) {
	f, err := os.Create(time.Now().Format("2006-01-02-15-04-05") + ".h264")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createWebsocket(name string, websocketUrl string) (*websocket.Conn, error) {
	log.Printf("CreateZoomSharingStream: dialing url= %v", websocketUrl)
	dialer := websocket.Dialer{
		EnableCompression: true,
	}

	websocketHeaders := http.Header{}
	websocketHeaders.Set("Accept-Language", "en-US,en;q=0.9")
	websocketHeaders.Set("Cache-Control", "no-cache")
	websocketHeaders.Set("Origin", "https://zoom.us")
	websocketHeaders.Set("Pragma", "no-cache")
	websocketHeaders.Set("User-Agent", userAgent)

	connection, _, err := dialer.Dial(websocketUrl, websocketHeaders)
	if err != nil {
		return nil, err
	}

	log.Printf("Dialed : %v", websocketUrl)

	return connection, nil
}

func (sharingStream *ZoomSharingStream) StartReceiveChannel() {
	connection := sharingStream.recv

	closeHandler := func(i int, msg string) error {
		log.Printf("Closing : %v %v", i, msg)
		return nil
	}
	connection.SetCloseHandler(closeHandler)

	recorder, err := Recorder()
	if err != nil {
		log.Fatal(err)
	}
	defer recorder.Close()

	decoder := sharingStream.decoder

	for {
		messageType, p, err := connection.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}

		// Pong
		if p[0] == PING {
			err := connection.WriteMessage(websocket.BinaryMessage, p)
			if err != nil {
				log.Fatal(err)
				return
			}
			// RTP packet
		} else if p[0] == RTP {
			log.Printf("pkt = %v", hex.EncodeToString(p))
			frame, err := decoder.Decode(p[4:])
			if err != nil {
				log.Fatal(err)
				return
			}
			_, err = recorder.Write(frame)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else if p[0] == RTCP {
			_, err := RtcpProcess(p[4:])
			if err != nil {
				log.Fatal(err)
				return
			}
		} else if p[0] == AES_GCM_IV_VALUE {
			// log.Printf("AES_GCM_IV_VALUE IV=%v", p[4:])
		} else {
			log.Printf("name=receive type=%v len=%v payload=%v", messageType, len(p), p)
		}
	}

}

func (sharingStream *ZoomSharingStream) SetSharedMeetingKey(encryptionKey string) error {
	sharedMeetingKey, err := ZoomEscapedBase64Decode(encryptionKey)
	if err != nil {
		return err
	}

	sharingStream.decoder.SetSharedMeetingKey(sharedMeetingKey)
	return nil
}
