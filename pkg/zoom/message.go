package zoom

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/gorilla/websocket"
)

var msgTypes = map[int]reflect.Type{
	WS_CONN_KEEPALIVE:       reflect.TypeOf(WebsocketConnectionKeepalive{}),
	WS_CONF_JOIN_RES:        reflect.TypeOf(JoinConferenceResponse{}),
	WS_CONF_CHAT_INDICATION: reflect.TypeOf(ConferenceChatIndication{}),
	// sender implemented, working
	WS_CONF_CHAT_REQ:             reflect.TypeOf(ConferenceChatRequest{}),
	WS_CONF_ATTRIBUTE_INDICATION: reflect.TypeOf(ConferenceAttributeIndication{}),
	WS_CONF_ROSTER_INDICATION:    reflect.TypeOf(ConferenceRosterIndication{}),
	// sender implemented, working
	WS_AUDIO_VOIP_JOIN_CHANNEL_REQ: reflect.TypeOf(AudioVoipJoinChannelRequest{}),
	// sender implemented, working
	WS_AUDIO_MUTE_REQ: reflect.TypeOf(AudioMuteRequest{}),
	// sender implemented, working
	WS_CONF_SET_SHARE_STATUS_REQ: reflect.TypeOf(SetShareStatusRequest{}),
	// sender implemented, working
	WS_CONF_RENAME_REQ: reflect.TypeOf(ConferenceRenameRequest{}),
	// sender implemented, untested
	WS_AUDIO_MUTEALL_REQ: reflect.TypeOf(AudioMuteAllRequest{}),
	// sender implemented, untested
	WS_CONF_SET_MUTE_UPON_ENTRY_REQ: reflect.TypeOf(ConferenceSetMuteUponEntryRequest{}),
	// sender implemented, untested
	WS_CONF_ALLOW_UNMUTE_AUDIO_REQ: reflect.TypeOf(ConferenceAllowUnmuteAudioRequest{}),
	// sender implemented, untested
	WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ: reflect.TypeOf(ConferenceAllowParticipantRenameRequest{}),
	// sender implemented, untested
	WS_CONF_ALLOW_UNMUTE_VIDEO_REQ: reflect.TypeOf(ConferenceAllowUnmuteVideoRequest{}),
	// sender implemented, working
	// yes there's a typo here, tell that to zoom
	WS_CONF_CHAT_PRIVILEDGE_REQ:       reflect.TypeOf(ConferenceChatPrivilegeRequest{}),
	WS_CONF_AVATAR_PERMISSION_CHANGED: reflect.TypeOf(ConferenceAvatarPermissionChanged{}),
	// sender implemented, untested
	WS_CONF_LOCK_SHARE_REQ:           reflect.TypeOf(ConferenceLockShareRequest{}),
	WS_CONF_LOCAL_RECORD_INDICATION:  reflect.TypeOf(ConferenceLocalRecordIndication{}),
	WS_CONF_OPTION_INDICATION:        reflect.TypeOf(ConferenceOptionIndication{}),
	WS_CONF_DC_REGION_INDICATION:     reflect.TypeOf(ConferenceDCRegionIndication{}),
	WS_AUDIO_SSRC_INDICATION:         reflect.TypeOf(SSRCIndication{}),
	WS_VIDEO_SSRC_INDICATION:         reflect.TypeOf(SSRCIndication{}),
	WS_VIDEO_ACTIVE_INDICATION:       reflect.TypeOf(VideoActiveIndication{}),
	WS_SHARING_STATUS_INDICATION:     reflect.TypeOf(SharingStatusIndication{}),
	WS_CONF_COHOST_CHANGE_INDICATION: reflect.TypeOf(ConferenceCohostChangeIndication{}),
	WS_AUDIO_ASN_INDICATION:          reflect.TypeOf(AudioAsnIndication{}),
	WS_CONF_BO_ATTRIBUTE_INDICATION:  reflect.TypeOf(ConferenceBreakoutRoomAttributeIndication{}),
	WS_CONF_BO_COMMAND_INDICATION:    reflect.TypeOf(ConferenceBreakoutRoomCommandIndication{}),
	// sender implemented, untested
	WS_CONF_BO_START_REQ: reflect.TypeOf(ConferenceBreakoutRoomStartRequest{}),
	// sender implemented, untested
	WS_CONF_BO_BROADCAST_REQ: reflect.TypeOf(ConferenceBreakoutRoomBroadcastRequest{}),
	// sender implemented, working
	WS_CONF_BO_JOIN_REQ: reflect.TypeOf(ConferenceBreakoutRoomJoinRequest{}),
	WS_CONF_BO_JOIN_RES: reflect.TypeOf(ConferenceBreakoutRoomJoinResponse{}),
	// sender implemented, untested
	WS_CONF_END_REQ: reflect.TypeOf(ConferenceEndRequest{}),
	// sender implemented, doesn't work???
	WS_CONF_BO_TOKEN_BATCH_REQ:     reflect.TypeOf(ConferenceBreakoutRoomTokenBatchRequest{}),
	WS_CONF_BO_TOKEN_RES:           reflect.TypeOf(ConferenceBreakoutRoomTokenResponse{}),
	WS_CONF_HOST_CHANGE_INDICATION: reflect.TypeOf(ConferenceHostChangeIndication{}),
	WS_CONF_END_INDICATION:         reflect.TypeOf(ConferenceEndIndication{}),
}

func GetMessageBody(message *GenericZoomMessage) (interface{}, error) {
	typ := msgTypes[message.Evt]
	if typ == nil {
		return nil, fmt.Errorf("Failed to get message body: missing type definition in zoom/message.go")
	}
	p := reflect.New(typ).Interface()
	err := json.Unmarshal(message.Body, p)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse body JSON: %+v", err)
	}
	return p, nil
}

func (session *ZoomSession) SendMessage(connection *websocket.Conn, eventNumber int, body interface{}) error {
	session.mutex.Lock() // gorilla/websocket only allows for 1 sender at a time + the send sequence number shouldn't be written to simultaneously
	defer session.mutex.Unlock()

	// increment number
	session.sendSequenceNumber += 1

	newMessage := GenericZoomMessage{
		Evt: eventNumber,
		Seq: session.sendSequenceNumber,
	}

	if body != nil {
		// body is a json.rawmessage so we have to do this
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		newMessage.Body = bodyBytes
	}
	log.Printf("Sending message (Evt: %s; Seq: %d): %s", MessageNumberToName[newMessage.Evt], newMessage.Seq, string(newMessage.Body))

	return connection.WriteJSON(newMessage)
}
