package zoom

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

/*
WEBSOCKET MESSAGE TYPES
*/

type Message interface {
}

type GenericZoomMessage struct {
	Body json.RawMessage `json:"body,omitempty"`
	Evt  int             `json:"evt"`
	Seq  uint32          `json:"seq"` // only positive - need this for atomic incrementer
}

type ConferenceChatIndication struct {
	AttendeeNodeID int                  `json:"attendeeNodeID"`
	DestNodeID     int                  `json:"destNodeID"`
	SenderName     BytesBase64NoPadding `json:"senderName"`
	Text           BytesBase64NoPadding `json:"text"`
}

type ConferenceChatRequest struct {
	DestNodeID int                  `json:"destNodeID"`
	Sn         BytesBase64NoPadding `json:"sn"`
	Text       BytesBase64NoPadding `json:"text"`
}

type JoinConferenceResponse struct {
	ConID         string `json:"conID"`
	ConfID        string `json:"confID"`
	E2EEncrypt    bool   `json:"e2eEncrypt"`
	Elapsed       int    `json:"elapsed"`
	EncType       int    `json:"encType"`
	HugeBO        bool   `json:"hugeBO"`
	Mn            string `json:"mn"`
	ParticipantID int    `json:"participantID"`
	Res           int    `json:"res"`
	Role          int    `json:"role"`
	SvcURL        string `json:"svcUrl"`
	UserGUID      string `json:"userGUID"`
	UserID        int    `json:"userID"`
	ZoomID        string `json:"zoomID"`
}

type WebsocketConnectionKeepalive struct{}
type ConferenceAttributeIndication map[string]interface{}

// there are many types of roster indication messages so we just omitempty everything so that we aren't sending a bunch of blank strings etc
type ConferenceRosterIndication struct {
	Add []struct {
		Avatar             string               `json:"avatar,omitempty"`
		BCCEditor          bool                 `json:"bCCEditor,omitempty"`
		BCanPinMultiVideo  bool                 `json:"bCanPinMultiVideo,omitempty"`
		BCapsPinMultiVideo bool                 `json:"bCapsPinMultiVideo,omitempty"`
		BGuest             bool                 `json:"bGuest,omitempty"`
		BHold              bool                 `json:"bHold,omitempty"`
		BRaiseHand         bool                 `json:"bRaiseHand,omitempty"`
		Dn2                BytesBase64NoPadding `json:"dn2,omitempty"`
		ID                 int                  `json:"id,omitempty"`
		Os                 int                  `json:"os,omitempty"`
		Role               int                  `json:"role,omitempty"`
		Type               int                  `json:"type,omitempty"`
		ZoomID             string               `json:"zoomID,omitempty"`
	} `json:"add"`
	Update []struct {
		// all these fields are optional
		Caps     int                  `json:"caps,omitempty"`
		Dn2      BytesBase64NoPadding `json:"dn2,omitempty"` // renames
		ID       int                  `json:"id,omitempty"`
		Muted    bool                 `json:"muted,omitempty"`
		BVideoOn bool                 `json:"bVideoOn,omitempty"`
		Audio    string               `json:"audio,omitempty`
		BCoHost  bool                 `json:"bCoHost,omitempty"`
		Role     int                  `json:"role,omitempty"`
	} `json:"update"`
	Remove []struct {
		ID          int `json:"id,omitempty"`
		NUSerStatus int `json:"nUserStatus,omitempty"`
	} `json:"remove"`
}

type ConferenceEndIndication struct {
	Reason    int `json:"reason"`
	SubReason int `json:"subReason"`
}

type AudioMuteRequest struct {
	BMute bool `json:"bMute"`
	ID    int  `json:"id,omitempty"` // default is you, if you specify id (I think??) you can get other people
}

type ConferenceRenameRequest struct {
	ID     int                  `json:"id"`
	Dn2    BytesBase64NoPadding `json:"dn2"`
	Olddn2 BytesBase64NoPadding `json:"olddn2"`
}

type AudioMuteAllRequest struct {
	BMute bool `json:"bMute"`
}

// TODO: merge these
type BOnRequest struct {
	BOn bool `json:"bOn"`
	ID  int  `json:"id,omitempty"` // default is you, if you specify id (I think??) you can get other people
}

type ConferenceSetShareStatusRequest struct {
	BOnRequest
	BShareAudio bool `json:"bShareAudio"`
}
type ConferenceSetMuteUponEntryRequest BOnRequest
type ConferenceAllowUnmuteAudioRequest BOnRequest
type ConferenceAllowParticipantRenameRequest BOnRequest
type ConferenceAllowUnmuteVideoRequest BOnRequest
type AudioVoipJoinChannelRequest BOnRequest
type VideoMuteRequest BOnRequest

type AudioVoipStatusRequest struct {
	OldAudioConnectionStatus int `json:"oldAudioConnectionStatus"`
	AudioConnectionStatus    int `json:"audioConnectionStatus"`
}

type ConferenceChatPrivilegeRequest struct {
	ChatPriviledge int `json:"chatPriviledge"`
}

type ConferenceAvatarPermissionChanged struct {
	BAllowedAvatar bool `json:"bAllowedAvatar"`
}

type ConferenceLockShareRequest struct {
	LockShare int `json:"lockShare"`
}

type ConferenceEndRequest struct{}

type ConferenceLocalRecordIndication struct{}

type ConferenceOptionIndication struct {
	DestNodeID int `json:"destNodeID"`
	// although this could be represented as bytesbase64nopadding data there is no point in doing so because we are not going to be anything other than copying it verbatim
	Opt string `json:"opt"`
	// Opt        BytesBase64NoPadding `json:"opt"`
}

// which datacenter are you using
type ConferenceDCRegionIndication struct {
	DC      string `json:"dc"`
	Network string `json:"network"`
	Region  string `json:"region"`
}

type SSRCIndication struct {
	SSRC int `json:"ssrc"`
}

type VideoActiveIndication struct {
	BVideoOn bool `json:"bVideoOn"`
	ID       int  `json:"id"`
}

type ConferenceHostChangeIndication struct {
	BCoHost bool `json:"bCoHost"`
}

type ConferenceCohostChangeIndication struct {
	BCoHost bool `json:"bCoHost"`
}

type AudioAsnIndication struct {
	Asn1 int `json:"asn1"`
}

type SharingStatusIndication struct {
	ActiveNodeID int `json:"activeNodeID"`
	BStatus      int `json:"bStatus"`
	Ssrc         int `json:"ssrc"`
}

type ConferenceBreakoutRoomTokenBatchRequest struct {
	Topic string `json:"topic"`
	Index int    `json:"index"`
}

type ConferenceBreakoutRoomTokenResponse struct {
	Bid string `json:"bid"`
}

type ConferenceBreakoutRoomBroadcastRequest struct {
	TextContent BytesBase64NoPadding `json:"textContent"`
}

type ConferenceBreakoutRoomJoinRequest struct {
	TargetBID string `json:"targetBID"`
}

type ConferenceBreakoutRoomJoinResponse struct {
	Bid     string `json:"bid"`
	Botoken string `json:"botoken"`
	ConfID  string `json:"confID"`
}

type ConferenceBreakoutRoomCommandIndication struct {
	// although this could be represented as bytesbase64nopadding data there is no point in doing so because we are not going to be anything other than copying it verbatim
	Botoken     string               `json:"botoken,omitempty"`
	CommandType string               `json:"commandType,omitempty"`
	ConfID      string               `json:"confID,omitempty"`
	RequestID   int                  `json:"requestID,omitempty"`
	TargetBID   string               `json:"targetBID,omitempty"`
	TextContent BytesBase64NoPadding `json:"textContent,omitempty"`
}

type ConferenceBreakoutRoomAttributeIndication struct {
	Proto ConferenceBreakoutRoomAttributeIndicationDataAlias `json:"proto"`
}

// start request is the same thing
type ConferenceBreakoutRoomStartRequest ConferenceBreakoutRoomAttributeIndication

type ConferenceBreakoutRoomAttributeIndicationData struct {
	ControlStatus              int                `json:"ControlStatus"`
	NameIndex                  int                `json:"NameIndex"`
	IsAutoJoinEnabled          bool               `json:"IsAutoJoinEnabled"`
	IsBackToMainSessionEnabled bool               `json:"IsBackToMainSessionEnabled"`
	IsTimerEnabled             bool               `json:"IsTimerEnabled"`
	TimerDuration              int                `json:"TimerDuration"`
	IsTimerAutoEndEnabled      bool               `json:"IsTimerAutoEndEnabled"`
	WaitSeconds                int                `json:"WaitSeconds"`
	StartTimeOnMMR             interface{}        `json:"StartTimeOnMMR"`
	ItemList                   []BreakoutRoomItem `json:"ItemList"`
	// ItemList                   []struct {
	// 	BID             string   `json:"BID"`
	// 	MeetingTitle    string   `json:"MeetingTitle"`
	// 	MeetingToken    string   `json:"MeetingToken"`
	// 	Status          int      `json:"Status"`
	// 	HostID          string   `json:"HostID"`
	// 	ParticipantList []string `json:"ParticipantList"`
	// } `json:"ItemList"`
}

type BreakoutRoomItem struct {
	BID             string   `json:"BID"`
	MeetingTitle    string   `json:"MeetingTitle"`
	MeetingToken    string   `json:"MeetingToken"`
	Status          int      `json:"Status"`
	HostID          string   `json:"HostID"`
	ParticipantList []string `json:"ParticipantList"`
}

type ConferenceHoldChangeIndication struct {
	BHold bool `json:"bHold"`
}

type SharingEncryptKeyIndication struct {
	AdditionalType int    `json:"additionalType"`
	EncryptKey     string `json:"encryptKey"`
}

type SharingSubscribeRequest struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

type SharingAssignedSendingSsrcResponse struct {
	SSRC int `json:"ssrc"`
}

// WebRTC audio related
type SharingReceivingChannelReadyIndication struct {
	SSRC        int  `json:"ssrc"`
	StreamIndex int  `json:"streamIndex"`
	VideoMode   bool `json:"videoMode"`
}

type SharingReceivingChannelCloseIndication struct {
	SSRC int `json:"ssrc"`
}

type DataChannelSendOfferToRWG struct {
	Offer string `json:"offer"`
	Type  int    `json:"type"`
}

type ConferenceBreakoutRoomAttributeIndicationDataAlias ConferenceBreakoutRoomAttributeIndicationData

func (b *ConferenceBreakoutRoomAttributeIndicationDataAlias) UnmarshalJSON(data []byte) error {
	dataUnquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	dataBytes, err := base64.StdEncoding.DecodeString(dataUnquoted)
	if err != nil {
		return err
	}

	real := ConferenceBreakoutRoomAttributeIndicationData{}
	err = json.Unmarshal(dataBytes, &real)
	if err != nil {
		return err
	}
	*b = (ConferenceBreakoutRoomAttributeIndicationDataAlias(real))

	return err
}

func (b ConferenceBreakoutRoomAttributeIndicationDataAlias) MarshalJSON() ([]byte, error) {
	jsonBytes, err := json.Marshal(ConferenceBreakoutRoomAttributeIndicationData(b))
	if err != nil {
		return nil, err
	}
	return []byte("\"" + base64.StdEncoding.EncodeToString([]byte(jsonBytes)) + "\""), nil
}

// golang usually allows []bytes to be represented by b64 strings but zoom sends them with no padding which breaks it so this is a custom type to allow for lack of padding
type BytesBase64NoPadding []byte

func (b *BytesBase64NoPadding) UnmarshalJSON(data []byte) error {
	dataUnquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	dataBytes, err := base64.RawURLEncoding.DecodeString(dataUnquoted)
	if err != nil {
		return err
	}
	bytesString := BytesBase64NoPadding(dataBytes)

	*b = bytesString
	return nil
}

func (b BytesBase64NoPadding) MarshalJSON() ([]byte, error) {
	return []byte("\"" + base64.RawURLEncoding.EncodeToString([]byte(b)) + "\""), nil
}
