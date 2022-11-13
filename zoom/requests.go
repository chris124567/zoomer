package zoom

func (session *ZoomSession) SendChatMessage(destNodeID int, text string) error {
	sendBody := ConferenceChatRequest{
		DestNodeID: destNodeID,
		Sn:         []byte(session.JoinInfo.ZoomID),
		Text:       []byte(text),
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_CHAT_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil

}

// host required
func (session *ZoomSession) RequestBreakoutRoomToken(topic string, index int) error {
	sendBody := ConferenceBreakoutRoomTokenBatchRequest{
		Topic: topic,
		Index: index,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_BO_TOKEN_BATCH_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
// request room bIDs using session.RequestBreakoutRoomToken, store them somewhere, then use those to make the rooms.  see struct details in message_types.go
func (session *ZoomSession) CreateBreakoutRoom(rooms []BreakoutRoomItem, autoJoin bool, timerEnabled bool, timerDurationSeconds int, forceLeaveWait int) error {
	protoData := ConferenceBreakoutRoomAttributeIndicationData{
		ControlStatus:     2,
		NameIndex:         1,
		IsAutoJoinEnabled: autoJoin,
		IsTimerEnabled:    timerEnabled,
		TimerDuration:     timerDurationSeconds,
		// how long before people are forced to leave
		WaitSeconds: forceLeaveWait,
		// ???
		StartTimeOnMMR: 464,
		ItemList:       rooms,
	}

	sendBody := ConferenceBreakoutRoomStartRequest{
		Proto: ConferenceBreakoutRoomAttributeIndicationDataAlias(protoData),
	}

	err := session.SendMessage(session.websocketConnection, WS_CONF_BO_START_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) BreakoutRoomBroadcast(text string) error {
	sendBody := ConferenceBreakoutRoomBroadcastRequest{
		TextContent: []byte(text),
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_BO_BROADCAST_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

/*
breakout room joining is not fully implemented
you have to send the WS_CONF_BO_JOIN_REQ (which this function does), wait for the WS_CONF_BO_JOIN_RES, then make a separate websocket connection using the token you get
breakout rooms are basically meetings= within meetings
*/
func (session *ZoomSession) RequestBreakoutRoomJoinToken(targetBID string) error {
	sendBody := ConferenceBreakoutRoomJoinRequest{
		TargetBID: targetBID,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_BO_JOIN_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// equivalent to zoom "join audio" - basically just allows us to have the voice icon next to our name
func (session *ZoomSession) JoinAudioVoipChannel(status bool) error {
	sendBody := AudioVoipJoinChannelRequest{
		ID:  session.JoinInfo.UserID,
		BOn: status,
	}
	err := session.SendMessage(session.websocketConnection, WS_AUDIO_VOIP_JOIN_CHANNEL_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// Signal that our audio channel is ready
func (session *ZoomSession) SignalAudioStatus(oldAudioConnectionStatus, audioConnectionStatus int) error {
	sendBody := AudioVoipStatusRequest{
		OldAudioConnectionStatus: oldAudioConnectionStatus,
		AudioConnectionStatus:    audioConnectionStatus,
	}
	err := session.SendMessage(session.websocketConnection, WS_AUDIO_VOIP_JOIN_CHANNEL_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// {"evt":8203,"body":{"oldAudioConnectionStatus":0,"audioConnectionStatus":1},"seq":3}
// {"evt":8203,"body":{"oldAudioConnectionStatus":0,"audioConnectionStatus":3,"seq":3}

// NOTE: this does not actually allow you to screenshare, that has yet to implemented.  it just changes the indicator next to your name and will show that you have solid black video
// true for mute, false for unmute
func (session *ZoomSession) SetVideoMuted(status bool) error {
	sendBody := VideoMuteRequest{
		ID:  session.JoinInfo.UserID,
		BOn: status,
	}
	err := session.SendMessage(session.websocketConnection, WS_VIDEO_MUTE_VIDEO_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// NOTE: this does not actually allow you to screenshare, that has yet to implemented.  it will show that you are sharing your screen but the output will be black
// true for mute, false for unmute
func (session *ZoomSession) SetShareStatus(status bool, shareAudio bool) error {
	sendBody := ConferenceSetShareStatusRequest{
		BOnRequest: BOnRequest{
			ID:  session.JoinInfo.UserID,
			BOn: status,
		},
		BShareAudio: shareAudio,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_SET_SHARE_STATUS_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// NOTE: this does not actually allow you to output any audio, that has yet to implemented.  it just changes the indicator next to your username.
// true for mute, false for unmute
func (session *ZoomSession) SetAudioMuted(status bool) error {
	// need to be into voip channel to be able to mute ourselves
	err := session.JoinAudioVoipChannel(true)
	if err != nil {
		return err
	}

	sendBody := AudioMuteRequest{
		BMute: status,
	}
	err = session.SendMessage(session.websocketConnection, WS_AUDIO_MUTE_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

func (session *ZoomSession) RenameMe(newName string) error {
	return session.RenameById(session.JoinInfo.UserID, session.Username, newName)
}

// host required to rename others (not self)
func (session *ZoomSession) RenameById(id int, oldName string, newName string) error {
	sendBody := ConferenceRenameRequest{
		ID:     id,
		Dn2:    []byte(newName),
		Olddn2: []byte(oldName),
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_RENAME_REQ, sendBody)
	if err != nil {
		return err
	}
	if id == session.JoinInfo.UserID {
		session.Username = newName
	}
	return nil
}

// host required
func (session *ZoomSession) RequestAllMute() error {
	sendBody := AudioMuteAllRequest{
		BMute: true,
	}
	err := session.SendMessage(session.websocketConnection, WS_AUDIO_MUTEALL_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) SetMuteUponEntry(status bool) error {
	sendBody := ConferenceSetMuteUponEntryRequest{
		BOn: status,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_SET_MUTE_UPON_ENTRY_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) SetAllowUnmuteAudio(status bool) error {
	sendBody := ConferenceAllowUnmuteAudioRequest{
		BOn: true,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_UNMUTE_AUDIO_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) SetAllowParticipantRename(status bool) error {
	sendBody := ConferenceAllowParticipantRenameRequest{
		BOn: true,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) SetAllowUnmuteVideo(status bool) error {
	sendBody := ConferenceAllowUnmuteVideoRequest{
		BOn: true,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_UNMUTE_VIDEO_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
// possible values: CHAT_EVERYONE_PUBLICLY_PRIVATELY = 1, CHAT_HOST_ONLY = 3, CHAT_NO_ONE = 4, CHAT_EVERYONE_PUBLICLY = 5
func (session *ZoomSession) SetChatLevel(status int) error {
	sendBody := ConferenceChatPrivilegeRequest{
		ChatPriviledge: status,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_CHAT_PRIVILEDGE_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

/*
host required
possible values:
CMM_SHARE_SETTING_HOST_GRAB = 0 (How many participants can share at the same time? One participant can share at a time) (Who can share? All Participants) (Who can start sharing when someone else is sharing? Only Host)

CMM_SHARE_SETTING_LOCK_SHARE  = 1 (How many participants can share at the same time? One participant can share at a time) (Who can share? Only Host)

CMM_SHARE_SETTING_ANYONE_GRAB = 2 (How many participants can share at the same time? One participant can share at a time) (Who can share? All Participants) (Who can start sharing when someone else is sharing? Only Host)

CMM_SHARE_SETTING_MULTI_SHARE = 3 (How many participants can share at the same time? Multiple participants can share simultaneously)
*/
func (session *ZoomSession) SetShareLockedStatus(status int) error {
	sendBody := ConferenceLockShareRequest{
		LockShare: status,
	}
	err := session.SendMessage(session.websocketConnection, WS_CONF_LOCK_SHARE_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

func (session *ZoomSession) SharingSubscribeRequest(id int, size int) error {
	sendBody := SharingSubscribeRequest{
		ID:   id,
		Size: size,
	}
	err := session.SendMessage(session.websocketConnection, WS_SHARING_SUBSCRIBE_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}

// host required
func (session *ZoomSession) EndMeeting() error {
	sendBody := ConferenceEndRequest{}
	err := session.SendMessage(session.websocketConnection, WS_CONF_END_REQ, sendBody)
	if err != nil {
		return err
	}
	return nil
}
