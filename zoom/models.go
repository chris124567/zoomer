package zoom

import (
	"encoding/json"
)

type MeetingInfo struct {
	Status    bool `json:"status"`
	ErrorCode int  `json:"errorCode"`
	Result    struct {
		Password                string                   `json:"passWord"`
		Auth                    string                   `json:"auth"`
		IsWebinar               int                      `json:"isWebinar,string"`
		TrackAuth               string                   `json:"track_auth"`
		Sign                    string                   `json:"sign"`
		EncryptedRWC            EncryptedRWCServersAlias `json:"encryptedRWC"`
		Mid                     string                   `json:"mid"`
		MeetingTopic            string                   `json:"meetingTopic"`
		Tid                     string                   `json:"tid"`
		CcURL                   string                   `json:"ccUrl"`
		TollNumbersJSON         string                   `json:"toll_numbers_json"`
		MeetingOptions          MeetingOptionsAlias      `json:"meetingOptions"`
		IsSupportQA             bool                     `json:"isSupportQA,string"`
		CallOutCountryJSON      CallOutCountryAlias      `json:"call_out_country_json"`
		H323Password            string                   `json:"h323Password"`
		RwcAgentEndpoint        string                   `json:"rwc_agent_endpoint"`
		UserName                string                   `json:"userName"`
		InviteEmail             string                   `json:"invite_email"`
		IsOriginHost            bool                     `json:"isOriginHost,string"`
		MeetingNumber           string                   `json:"meetingNumber"`
		FormateConfno           string                   `json:"formate_confno"`
		OptionVideoHost         bool                     `json:"optionVideoHost,string"`
		SupportCallOut          int                      `json:"support_call_out,string"`
		OptionAudioType         string                   `json:"optionAudioType"`
		RwcAgentEndpointBackup  string                   `json:"rwc_agent_endpoint_backup"`
		Location                string                   `json:"location"`
		CanRecording            int                      `json:"can_recording,string"`
		OptionVideoParticipants string                   `json:"optionVideoParticipants"`
		Ts                      int64                    `json:"ts,string"`
	} `json:"result"`
}

type MeetingOptions struct {
	IsShareWhiteboardEnabled                bool `json:"isShareWhiteboardEnabled"`
	IsChatEnabled                           bool `json:"isChatEnabled"`
	IsAllowBreakoutRoomPreAssign            bool `json:"isAllowBreakoutRoomPreAssign"`
	AllowParticipantsRename                 bool `json:"allowParticipantsRename"`
	IsShareOptionLocked                     bool `json:"isShareOptionLocked"`
	AllowPutOnHold                          bool `json:"allowPutOnHold"`
	IsEnableHideParticipantPic              bool `json:"isEnableHideParticipantPic"`
	AllowParticipantsRenameLocked           bool `json:"allowParticipantsRenameLocked"`
	IsScreenShareEnabled                    bool `json:"isScreenShareEnabled"`
	IsBO100Enabled                          bool `json:"isBO100Enabled"`
	IsPollingEnabled                        bool `json:"isPollingEnabled"`
	IsEnableMuteParticipantsUponEntry       bool `json:"isEnableMuteParticipantsUponEntry"`
	IsGroupHDEnabled                        bool `json:"isGroupHDEnabled"`
	IsEnableMeetingWatermark                bool `json:"isEnableMeetingWatermark"`
	IsAllowShareZMWindowEnabled             bool `json:"isAllowShareZMWindowEnabled"`
	IsReportParticipantsEnabled             bool `json:"isReportParticipantsEnabled"`
	IsEnableDeleteChatMsg                   bool `json:"isEnableDeleteChatMsg"`
	IsPrivateChatEnabled                    bool `json:"isPrivateChatEnabled"`
	IsUserEnableRecordingReminder           bool `json:"isUserEnableRecordingReminder"`
	IsLocalRecordingEnabled                 bool `json:"isLocalRecordingEnabled"`
	IsEnableClosedCaption                   bool `json:"isEnableClosedCaption"`
	IsAllowParticipantsReport               bool `json:"isAllowParticipantsReport"`
	IsAuthLocalRecording                    bool `json:"isAuthLocalRecording"`
	IsRemoteControlEnbaled                  bool `json:"isRemoteControlEnbaled"`
	IsInstantMeeting                        bool `json:"isInstantMeeting"`
	IsPrivateChatLocked                     bool `json:"isPrivateChatLocked"`
	IsEnableAutomaticDisplayJoinAudioDialog bool `json:"isEnableAutomaticDisplayJoinAudioDialog"`
	IsAllowChatAndRaiseChime                bool `json:"isAllowChatAndRaiseChime"`
	EnableWaitingRoom                       bool `json:"enableWaitingRoom"`
	IsMeetingScreenSharingGrabAll           bool `json:"isMeetingScreenSharingGrabAll"`
	IsChatLocked                            bool `json:"isChatLocked"`
	IsEnableLiveTranscription               bool `json:"isEnableLiveTranscription"`
	IsEnableMeetingControlToolBar           bool `json:"isEnableMeetingControlToolBar"`
	IsCOHostEnabled                         bool `json:"isCOHostEnabled"`
	IsEnableEncryption3RdParty              bool `json:"isEnableEncryption3rdParty"`
	IsShareScreenHostOnly                   bool `json:"isShareScreenHostOnly"`
	IsEnableBreakoutRoom                    bool `json:"isEnableBreakoutRoom"`
	PlayUserJoinLeaveAudio                  bool `json:"playUserJoinLeaveAudio"`
	Nonverbalfeedback                       bool `json:"nonverbalfeedback"`
	IsWaitingRoomLocked                     bool `json:"isWaitingRoomLocked"`
	IsWaterMarkLocked                       bool `json:"isWaterMarkLocked"`
}

type CallOutCountry []struct {
	Code string `json:"code"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type RwgInfo struct {
	Rwg     string `json:"rwg"`
	RwcAuth string `json:"rwcAuth"`
}

/*
axios does this annoying thing where it puts json within json as a string
golang has the "string" option to resolve this but this only works for numbers and booleans, so we have to make custom marshalers to solve this

this will produce inconsistencies when marshaling because i did not implement a marshaler that would quote the output and make it a string but since we are only receiving this data only unmarshalling really matters
*/

type MeetingOptionsAlias MeetingOptions

func (m *MeetingOptionsAlias) UnmarshalJSON(data []byte) error {
	// Try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		var originalTypeData MeetingOptions
		if err := json.Unmarshal([]byte(str), &originalTypeData); err != nil {
			return err
		}
		*m = MeetingOptionsAlias(originalTypeData)
		return nil
	}

	return json.Unmarshal(data, m)
}

type EncryptedRWCServersAlias map[string]string

func (m *EncryptedRWCServersAlias) UnmarshalJSON(data []byte) error {
	// Try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		var originalTypeData map[string]string
		if err := json.Unmarshal([]byte(str), &originalTypeData); err != nil {
			return err
		}
		*m = EncryptedRWCServersAlias(originalTypeData)
		return nil
	}

	return json.Unmarshal(data, m)
}

type CallOutCountryAlias CallOutCountry

func (m *CallOutCountryAlias) UnmarshalJSON(data []byte) error {
	// Try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		var originalTypeData CallOutCountry
		if err := json.Unmarshal([]byte(str), &originalTypeData); err != nil {
			return err
		}
		*m = CallOutCountryAlias(originalTypeData)
		return nil
	}

	return json.Unmarshal(data, m)
}
