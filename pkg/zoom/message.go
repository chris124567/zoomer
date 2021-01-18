package zoom

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func GetMessageBody(message *GenericZoomMessage) (interface{}, error) {
	// keep alive messages have no body
	if message.Evt == WS_CONN_KEEPALIVE {
		return nil, fmt.Errorf("Failed to get message body: Keep alive messages have no body")
	}

	bodyDataPointer := getPointerForBody(message)
	if bodyDataPointer == nil {
		return nil, fmt.Errorf("Failed to get message body: missing type definition in zoom/message.go")
	}

	err := json.Unmarshal(message.Body, &bodyDataPointer)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse body JSON: %+v", err)
	}
	return bodyDataPointer, nil

}

func getPointerForBody(message *GenericZoomMessage) interface{} {
	switch message.Evt {
	case WS_CONN_KEEPALIVE:
		var p WebsocketConnectionKeepalive
		return &p

	case WS_CONF_JOIN_RES:
		var p JoinConferenceResponse
		return &p

	case WS_CONF_CHAT_INDICATION:
		var p ConferenceChatIndication
		return &p

	// sender implemented, working
	case WS_CONF_CHAT_REQ:
		var p ConferenceChatRequest
		return &p

	case WS_CONF_ATTRIBUTE_INDICATION:
		var p ConferenceAttributeIndication
		return &p

	case WS_CONF_ROSTER_INDICATION:
		var p ConferenceRosterIndication
		return &p

	// sender implemented, working
	case WS_AUDIO_VOIP_JOIN_CHANNEL_REQ:
		var p AudioVoipJoinChannelRequest
		return &p

	// sender implemented, working
	case WS_AUDIO_MUTE_REQ:
		var p AudioMuteRequest
		return &p

	// sender implemented, working
	case WS_CONF_SET_SHARE_STATUS_REQ:
		var p SetShareStatusRequest
		return &p

	// sender implemented, working
	case WS_CONF_RENAME_REQ:
		var p ConferenceRenameRequest
		return &p

	// sender implemented, untested
	case WS_AUDIO_MUTEALL_REQ:
		var p AudioMuteAllRequest
		return &p

	// sender implemented, untested
	case WS_CONF_SET_MUTE_UPON_ENTRY_REQ:
		var p ConferenceSetMuteUponEntryRequest
		return &p

	// sender implemented, untested
	case WS_CONF_ALLOW_UNMUTE_AUDIO_REQ:
		var p ConferenceAllowUnmuteAudioRequest
		return &p

	// sender implemented, untested
	case WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ:
		var p ConferenceAllowParticipantRenameRequest
		return &p

	// sender implemented, untested
	case WS_CONF_ALLOW_UNMUTE_VIDEO_REQ:
		var p ConferenceAllowUnmuteVideoRequest
		return &p

	// sender implemented, working
	case WS_CONF_CHAT_PRIVILEDGE_REQ: // yes there's a typo here, tell that to zoom
		var p ConferenceChatPrivilegeRequest
		return &p

	case WS_CONF_AVATAR_PERMISSION_CHANGED:
		var p ConferenceAvatarPermissionChanged
		return &p

	// sender implemented, untested
	case WS_CONF_LOCK_SHARE_REQ:
		var p ConferenceLockShareRequest
		return &p

	case WS_CONF_LOCAL_RECORD_INDICATION:
		var p ConferenceLocalRecordIndication
		return &p

	case WS_CONF_OPTION_INDICATION:
		var p ConferenceOptionIndication
		return &p

	case WS_CONF_DC_REGION_INDICATION:
		var p ConferenceDCRegionIndication
		return &p

	case WS_AUDIO_SSRC_INDICATION:
		var p SSRCIndication
		return &p

	case WS_VIDEO_SSRC_INDICATION:
		var p SSRCIndication
		return &p

	case WS_VIDEO_ACTIVE_INDICATION:
		var p VideoActiveIndication
		return &p

	case WS_SHARING_STATUS_INDICATION:
		var p SharingStatusIndication
		return &p

	case WS_CONF_COHOST_CHANGE_INDICATION:
		var p ConferenceCohostChangeIndication
		return &p

	case WS_AUDIO_ASN_INDICATION:
		var p AudioAsnIndication
		return &p

	case WS_CONF_BO_ATTRIBUTE_INDICATION:
		var p ConferenceBreakoutRoomAttributeIndication
		return &p

	case WS_CONF_BO_COMMAND_INDICATION:
		var p ConferenceBreakoutRoomCommandIndication
		return &p

	// sender implemented, untested
	case WS_CONF_BO_START_REQ:
		var p ConferenceBreakoutRoomStartRequest
		return &p

	// sender implemented, untested
	case WS_CONF_BO_BROADCAST_REQ:
		var p ConferenceBreakoutRoomBroadcastRequest
		return &p

	// sender implemented, working
	case WS_CONF_BO_JOIN_REQ:
		var p ConferenceBreakoutRoomJoinRequest
		return &p

	case WS_CONF_BO_JOIN_RES:
		var p ConferenceBreakoutRoomJoinResponse
		return &p

	// sender implemented, untested
	case WS_CONF_END_REQ:
		var p ConferenceEndRequest
		return &p

	// sender implemented, doesn't work???
	case WS_CONF_BO_TOKEN_BATCH_REQ:
		var p ConferenceBreakoutRoomTokenBatchRequest
		return &p

	case WS_CONF_BO_TOKEN_RES:
		var p ConferenceBreakoutRoomTokenResponse
		return &p

	case WS_CONF_HOST_CHANGE_INDICATION:
		var p ConferenceHostChangeIndication
		return &p

	case WS_CONF_END_INDICATION:
		var p ConferenceEndIndication
		return &p
	}

	return nil
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
