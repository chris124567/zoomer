package zoom

func (session *ZoomSession) SendChatMessage(destNodeID int, text string) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_CHAT_REQ, ConferenceChatRequest{
		DestNodeID: destNodeID,
		Sn:         []byte(session.JoinInfo.ZoomID),
		Text:       []byte(text),
	})
}

// host required
func (session *ZoomSession) RequestBreakoutRoomToken(topic string, index int) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_BO_TOKEN_BATCH_REQ, ConferenceBreakoutRoomTokenBatchRequest{
		Topic: topic,
		Index: index,
	})
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
	return session.SendMessage(session.websocketConnection, WS_CONF_BO_START_REQ, ConferenceBreakoutRoomStartRequest{
		Proto: ConferenceBreakoutRoomAttributeIndicationDataAlias(protoData),
	})
}

// host required
func (session *ZoomSession) BreakoutRoomBroadcast(text string) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_BO_BROADCAST_REQ, ConferenceBreakoutRoomBroadcastRequest{
		TextContent: []byte(text),
	})
}

/*
breakout room joining is not fully implemented
you have to send the WS_CONF_BO_JOIN_REQ (which this function does), wait for the WS_CONF_BO_JOIN_RES, then make a separate websocket connection using the token you get
breakout rooms are basically meetings= within meetings
*/
func (session *ZoomSession) RequestBreakoutRoomJoinToken(targetBID string) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_BO_JOIN_REQ, ConferenceBreakoutRoomJoinRequest{
		TargetBID: targetBID,
	})
}

// equivalent to zoom "join audio" - basically just allows us to have the voice icon next to our name
func (session *ZoomSession) JoinAudioVoipChannel(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_AUDIO_VOIP_JOIN_CHANNEL_REQ, AudioVoipJoinChannelRequest{
		ID:  session.JoinInfo.UserID,
		BOn: status,
	})
}

// NOTE: this does not actually allow you to screenshare, that has yet to implemented.  it just changes the indicator next to your name and will show that you have solid black video
// true for mute, false for unmute
func (session *ZoomSession) SetVideoMuted(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_VIDEO_MUTE_VIDEO_REQ, VideoMuteRequest{
		ID:  session.JoinInfo.UserID,
		BOn: status,
	})
}

// NOTE: this does not actually allow you to screenshare, that has yet to implemented.  it will show that you are sharing your screen but the output will be black
// true for mute, false for unmute
func (session *ZoomSession) SetScreenShareMuted(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_SET_SHARE_STATUS_REQ, SetShareStatusRequest{
		ID:  session.JoinInfo.UserID,
		BOn: status,
	})
}

// NOTE: this does not actually allow you to output any audio, that has yet to implemented.  it just changes the indicator next to your username.
// true for mute, false for unmute
func (session *ZoomSession) SetAudioMuted(status bool) error {
	// need to be into voip channel to be able to mute ourselves
	if err := session.JoinAudioVoipChannel(true); err != nil {
		return err
	}
	return session.SendMessage(session.websocketConnection, WS_AUDIO_MUTE_REQ, AudioMuteRequest{
		BMute: status,
	})
}

func (session *ZoomSession) RenameMe(newName string) error {
	return session.RenameById(session.JoinInfo.UserID, session.Username, newName)
}

// host required to rename others (not self)
func (session *ZoomSession) RenameById(id int, oldName string, newName string) error {
	if err := session.SendMessage(session.websocketConnection, WS_CONF_RENAME_REQ, ConferenceRenameRequest{
		ID:     id,
		Dn2:    []byte(newName),
		Olddn2: []byte(oldName),
	}); err != nil {
		return err
	}
	if id == session.JoinInfo.UserID {
		session.Username = newName
	}
	return nil
}

// host required
func (session *ZoomSession) RequestAllMute() error {
	return session.SendMessage(session.websocketConnection, WS_AUDIO_MUTEALL_REQ, AudioMuteAllRequest{
		BMute: true,
	})
}

// host required
func (session *ZoomSession) SetMuteUponEntry(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_SET_MUTE_UPON_ENTRY_REQ, ConferenceSetMuteUponEntryRequest{
		BOn: status,
	})
}

// host required
func (session *ZoomSession) SetAllowUnmuteAudio(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_UNMUTE_AUDIO_REQ, ConferenceAllowUnmuteAudioRequest{
		BOn: true,
	})
}

// host required
func (session *ZoomSession) SetAllowParticipantRename(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ, ConferenceAllowParticipantRenameRequest{
		BOn: true,
	})
}

// host required
func (session *ZoomSession) SetAllowUnmuteVideo(status bool) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_ALLOW_UNMUTE_VIDEO_REQ, ConferenceAllowUnmuteVideoRequest{
		BOn: true,
	})
}

// host required
// possible values: CHAT_EVERYONE_PUBLICLY_PRIVATELY = 1, CHAT_HOST_ONLY = 3, CHAT_NO_ONE = 4, CHAT_EVERYONE_PUBLICLY = 5
func (session *ZoomSession) SetChatLevel(status int) error {
	return session.SendMessage(session.websocketConnection, WS_CONF_CHAT_PRIVILEDGE_REQ, ConferenceChatPrivilegeRequest{
		ChatPriviledge: status,
	})
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
	return session.SendMessage(session.websocketConnection, WS_CONF_LOCK_SHARE_REQ, ConferenceLockShareRequest{
		LockShare: status,
	})
}

// host required
func (session *ZoomSession) EndMeeting() error {
	return session.SendMessage(session.websocketConnection, WS_CONF_END_REQ, ConferenceEndRequest{})
}
