package zoom

import "net/http"

// make sure these are consistent
// TODO: make configurable
const (
	userAgent          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	userAgentShorthand = "Chrome112" // todo: figure out zooms algorithm for determining this
)

func httpHeaders() http.Header {
	return http.Header{
		http.CanonicalHeaderKey("pragma"):                    []string{"no-cache"},
		http.CanonicalHeaderKey("cache-control"):             []string{"no-cache"},
		http.CanonicalHeaderKey("upgrade-insecure-requests"): []string{"1"},
		http.CanonicalHeaderKey("user-agent"):                []string{userAgent},
		http.CanonicalHeaderKey("accept"):                    []string{"application/json, text/plain, */*"},
		http.CanonicalHeaderKey("sec-fetch-site"):            []string{"none"},
		http.CanonicalHeaderKey("sec-fetch-mode"):            []string{"navigate"},
		http.CanonicalHeaderKey("sec-fetch-user"):            []string{"?1"},
		http.CanonicalHeaderKey("sec-fetch-dest"):            []string{"document"},
		http.CanonicalHeaderKey("accept-language"):           []string{"en-US,en;q=0.9"},
	}
}

// from webclient.js
const (
	// other
	WS_CONN_KEEPALIVE     = 0 // WebsocketConnectionKeepalive
	NEED_UPDATE_WEBSDK    = 1
	CONF_EVT_TYPE_BASE    = 4096
	AUDIO_EVT_TYPE_BASE   = 8192
	VIDEO_EVT_TYPE_BASE   = 12288
	SHARING_EVT_TYPE_BASE = 16384
	XMPP_EVT_TYPE_BASE    = 24576
	// websocket
	WS_CONF_JOIN_REQ                                 = 4097
	WS_CONF_JOIN_RES                                 = 4098 // JoinConferenceResponse
	WS_CONF_LOCK_REQ                                 = 4099
	WS_CONF_LOCK_RES                                 = 4100
	WS_CONF_END_REQ                                  = 4101 // ConferenceEndRequest
	WS_CONF_END_RES                                  = 4102
	WS_CONF_LEAVE_REQ                                = 4103
	WS_CONF_LEAVE_RES                                = 4104
	WS_CONF_RECORD_REQ                               = 4105
	WS_CONF_RECORD_RES                               = 4106
	WS_CONF_EXPEL_REQ                                = 4107
	WS_CONF_EXPEL_RES                                = 4108
	WS_CONF_RENAME_REQ                               = 4109 // ConferenceRenameRequest
	WS_CONF_ASSIGN_HOST_REQ                          = 4111
	WS_CONF_PUT_ON_HOLD_REQ                          = 4113
	WS_CONF_SET_MUTE_UPON_ENTRY_REQ                  = 4115 // ConferenceSetMuteUponEntryRequest
	WS_CONF_SET_HOLD_UPON_ENTRY_REQ                  = 4117
	WS_CONF_INVITE_CRC_DEVICE_REQ                    = 4119
	WS_CONF_INVITE_CRC_DEVICE_RES                    = 4120
	WS_CONF_CANCEL_INVITE_CRC_DEVICE_REQ             = 4121
	WS_CONF_CANCEL_INVITE_CRC_DEVICE_RES             = 4122
	WS_CONF_SET_BROADCAST_REQ                        = 4123
	WS_CONF_SET_BROADCAST_RES                        = 4124
	WS_CONF_CLOSED_CAPTION_REQ                       = 4125
	WS_CONF_CLOSED_CAPTION_RES                       = 4126
	WS_CONF_ALLOW_VIEW_PARTICIPANT_REQ               = 4127
	WS_CONF_LOWER_ALL_HAND_REQ                       = 4129
	WS_CONF_RAISE_LOWER_HAND_REQ                     = 4131
	WS_CONF_RECLAIM_HOST_REQ                         = 4133
	WS_CONF_CHAT_REQ                                 = 4135 // ConferenceChatRequest
	WS_CONF_ASSIGN_CC_REQ                            = 4137
	WS_CONF_CHAT_PRIVILEDGE_REQ                      = 4141 // ConferenceChatPrivilegeRequest
	WS_CONF_FEEDBACK_REQ                             = 4143
	WS_CONF_FEEDBACK_CLEAR_REQ                       = 4145
	WS_CONF_ALLOW_UNMUTE_VIDEO_REQ                   = 4147 // ConferenceAllowUnmuteVideoRequest
	WS_CONF_ALLOW_UNMUTE_AUDIO_REQ                   = 4149 // ConferenceAllowUnmuteAudioRequest
	WS_CONF_ALLOW_RAISE_HAND_REQ                     = 4151
	WS_CONF_PANELIST_VOTE_REQ                        = 4153
	WS_CONF_BO_TOKEN_REQ                             = 4173
	WS_CONF_BO_TOKEN_RES                             = 4174 // ConferenceBreakoutRoomTokenResponse
	WS_CONF_BO_START_REQ                             = 4175 // ConferenceBreakoutRoomStartRequest
	WS_CONF_BO_STOP_REQ                              = 4177
	WS_CONF_BO_ASSIGN_REQ                            = 4179
	WS_CONF_BO_SWITCH_REQ                            = 4181
	WS_CONF_BO_WANT_JOIN_REQ                         = 4183
	WS_CONF_BO_LEAVE_REQ                             = 4185
	WS_CONF_BO_BROADCAST_REQ                         = 4187 // ConferenceBreakoutRoomBroadcastRequest
	WS_CONF_BO_HELP_REQ                              = 4189
	WS_CONF_BO_HELP_RESULT_REQ                       = 4191
	WS_CONF_BO_JOIN_REQ                              = 4193 // ConferenceBreakoutRoomJoinRequest
	WS_CONF_BO_JOIN_RES                              = 4194 // ConferenceBreakoutRoomJoinResponse
	WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ             = 4163 // ConferenceAllowParticipantRenameRequest
	WS_CONF_ALLOW_MESSAGE_FEEDBACK_NOTIFY_REQ        = 4171
	WS_CONF_REVOKE_COHOST_REQ                        = 4195
	WS_CONF_PLAY_CHIME_OPEN_CLOSE_REQ                = 4197
	WS_CONF_ADMIT_ALL_SILENT_USERS_REQ               = 4199
	WS_CONF_BIND_UNBIND_TELE_USR_REQ                 = 4201
	WS_CONF_ALLOW_QA_AUTO_REPLY_REQ                  = 4203
	WS_CONF_EXPEL_ATTENDEE_REQ                       = 4205
	WS_CONF_EXPEL_ATTENDEE_RES                       = 4206
	WS_CONF_PRACTICE_SESSION_REQ                     = 4207
	WS_CONF_PRACTICE_SESSION_RES                     = 4208
	WS_CONF_ROLE_CHANGE_REQ                          = 4209
	WS_CONF_ROLE_CHANGE_RES                          = 4210
	WS_CONF_BO_TOKEN_BATCH_REQ                       = 4211 // ConferenceBreakoutRoomTokenBatchRequest
	WS_CONF_BO_PRE_ASSIGN_REQ                        = 4213
	WS_CONF_BO_PRE_ASSIGN_RES                        = 4214
	WS_CONF_CHANGE_MULTI_PIN_PRIVILGE_REQ            = 4217
	WS_CONF_SET_GROUP_LAYOUT                         = 4219
	WS_CONF_HOST_KEY_REQ                             = 4215
	WS_CONF_HOST_KEY_RES                             = 4216
	WS_CONF_AVATAR_PERMISSION_CHANGED                = 4222 // ConferenceAvatarPermissionChanged
	WS_CONF_SUSPEND_MEETING                          = 4229
	WS_CONF_SUSPEND_MEETING_REQ_RESULT               = 4230
	WS_AUDIO_MUTE_REQ                                = 8193 // AudioMuteRequest
	WS_AUDIO_MUTE_RES                                = 8194
	WS_AUDIO_DROP_REQ                                = 8195
	WS_AUDIO_DROP_RES                                = 8196
	WS_AUDIO_DIALOUT_REQ                             = 8197
	WS_AUDIO_DIALOUT_RES                             = 8198
	WS_AUDIO_CANCEL_DIALOUT_REQ                      = 8199
	WS_AUDIO_CANCEL_DIALOUT_RES                      = 8200
	WS_AUDIO_MUTEALL_REQ                             = 8201 // AudioMuteAllRequest
	WS_AUDIO_MUTEALL_RES                             = 8202
	WS_AUDIO_ALLOW_TALK_REQ                          = 8204
	WS_AUDIO_ALLOW_TALK_RES                          = 8205
	WS_CONF_ROSTER_INDICATION                        = 7937 // ConferenceRosterIndication
	WS_CONF_ATTRIBUTE_INDICATION                     = 7938 // ConferenceAttributeIndication
	WS_CONF_END_INDICATION                           = 7939
	WS_CONF_HOST_CHANGE_INDICATION                   = 7940 // ConferenceHostChangeIndication
	WS_CONF_COHOST_CHANGE_INDICATION                 = 7941 // ConferenceCohostChangeIndication
	WS_CONF_HOLD_CHANGE_INDICATION                   = 7942 // ConferenceHoldChangeIndication
	WS_CONF_CLOSED_CAPTION_INDICATION                = 7943
	WS_CONF_CHAT_INDICATION                          = 7944 // ConferenceChatIndication
	WS_CONF_OPTION_INDICATION                        = 7945 // ConferenceOptionIndication
	WS_CONF_KV_UPDATE_INDICATION                     = 7946
	WS_CONF_LOCAL_RECORD_INDICATION                  = 7947 // ConferenceLocalRecordIndication
	WS_AUDIO_VOIP_JOIN_CHANNEL_REQ                   = 8203 // AudioVoipJoinChannelRequest
	WS_CONF_BO_COMMAND_INDICATION                    = 7949 // ConferenceBreakoutRoomCommandIndication
	WS_CONF_BO_ATTRIBUTE_INDICATION                  = 7950 // ConferenceBreakoutRoomAttributeIndication
	WS_CONF_ADMIT_ALL_SILENT_USERS_INDICATION        = 7951
	WS_CONF_BIND_UNBIND_INDICATION                   = 7952
	WS_CONF_UPDATE_MEETING_TOPIC_INDICATION          = 7953
	WS_CONF_DC_REGION_INDICATION                     = 7954 // ConferenceDCRegionIndication
	WS_CONF_CAN_ADMIT_WHEN_NOHOST_PRESENT_INDICATION = 7955
	WS_CONF_GROUP_LAYOUT_INDICATION                  = 7958
	WS_AUDIO_ASN_INDICATION                          = 12033 // AudioAsnIndication
	WS_AUDIO_MUTE_INDICATION                         = 12034
	WS_AUDIO_SSRC_INDICATION                         = 12035 // AudioSSRCIndication
	WS_AUDIO_ALLOW_TALK_INDICATION                   = 12036
	WS_AUDIO_SSRC_ASK_UNMUTE_INDICATION              = 12037
	WS_WEBINAR_VIEW_ONLY_TELEPHONY_INDICATION        = 12038
	WS_VIDEO_ACTIVE_INDICATION                       = 16129 // VideoActiveIndication
	WS_VIDEO_SSRC_INDICATION                         = 16131 // SSRCIndication
	WS_VIDEO_MUTE_INDICATION                         = 16133
	WS_VIDEO_LEADERSHIP_INDICATION                   = 16135
	WS_VIDEO_SUBSCRIBE_REQ                           = 12289
	WS_VIDEO_UNSUBSCRIBE_REQ                         = 12291
	WS_VIDEO_KEY_FRAME_REQ                           = 12293
	WS_VIDEO_NETWORK_FEEDBACK                        = 12295
	WS_VIDEO_MUTE_VIDEO_REQ                          = 12297
	WS_VIDEO_SPOTLIGHT_VIDEO_REQ                     = 12299
	WS_SHARING_PAUSE_REQ                             = 16385
	WS_SHARING_RESUME_REQ                            = 16387
	WS_SHARING_STATUS_INDICATION                     = 20225 // SharingStatusIndication
	WS_SHARING_SIZE_CHANGE_INDICATION                = 20226
	WS_CONF_ALLOW_ANONYMOUS_QUESTION_REQ             = 4155
	WS_CONF_ALLOW_VIEW_ALL_QUESTION_REQ              = 4157
	WS_CONF_ALLOW_UPVOTE_QUESTION_REQ                = 4159
	WS_CONF_ALLOW_COMMENT_QUESTION_REQ               = 4161
	WS_SHARING_REMOTE_CONTROL_REQ                    = 16389
	WS_SHARING_REMOTE_CONTROL_INDICATION             = 16391
	WS_SHARING_REMOTE_CONTROLLER_GRAB                = 16393
	WS_SHARING_REMOTE_CONTROLLER_GRAB_INDICATION     = 16395
	WS_SHARING_SUBSCRIBE_REQ                         = 16415
	WS_SHARING_UNSUBSCRIBE_REQ                       = 16417
	WS_MEETING_RWG_CONNECT_TIME                      = 4167
	EVT_TYPE_WS_VIDEO_DATACHANNEL_ANSWER             = 24322
	DATA_CHANNEL_SEND_OFFER_TO_RWG                   = 24321
	WS_CONF_FOLLOW_HOST_REQ                          = 4223
	WS_CONF_DRAG_LAYOUT_INDICATION                   = 7957
	WS_CONF_SET_DRAG_LAYOUT                          = 4218
	WS_CONF_LIVE_TRANSCRIPTION_ON_OFF_REQ            = 4227
	WS_CONF_LIVE_TRANSCRIPTION_ON_OFF_RES            = 4228
	WS_CONF_LIVE_TRANSCRIPTION_STATUS_INDICATION     = 7959
	WS_CONF_POLLING_REQ                              = 4165
	WS_CONF_POLLING_USER_ACTION_REQ                  = 4224
	WS_CONF_POLLING_USER_ACTION_ERROR                = 4225
	WS_CONF_POLLING_SET_POLLING_TOKEN                = 4226

	// OTHER NOT FOUND DIRECTLY IN JAVASCRIPT
	WS_CONF_LOCK_SHARE_REQ       = 4169  // ConferenceLockShareRequest
	WS_CONF_SET_SHARE_STATUS_REQ = 16409 // SetShareStatusRequest
)

// settings for WS_CONF_LOCK_SHARE_REQ
const (
	CMM_SHARE_SETTING_HOST_GRAB   = 0
	CMM_SHARE_SETTING_LOCK_SHARE  = 1
	CMM_SHARE_SETTING_ANYONE_GRAB = 2
	CMM_SHARE_SETTING_MULTI_SHARE = 3
)

const (
	CHAT_EVERYONE_PUBLICLY_PRIVATELY = 1
	CHAT_HOST_ONLY                   = 3
	CHAT_NO_ONE                      = 4
	CHAT_EVERYONE_PUBLICLY           = 5
)

const (
	EVERYONE_CHAT_ID = 0
)

// for debugging and logging purposes only - use enum in code
var MessageNumberToName map[int]string = map[int]string{
	0:     "WS_CONN_KEEPALIVE",
	1:     "NEED_UPDATE_WEBSDK",
	4096:  "CONF_EVT_TYPE_BASE",
	8192:  "AUDIO_EVT_TYPE_BASE",
	12288: "VIDEO_EVT_TYPE_BASE",
	16384: "SHARING_EVT_TYPE_BASE",
	24576: "XMPP_EVT_TYPE_BASE",
	4097:  "WS_CONF_JOIN_REQ",
	4098:  "WS_CONF_JOIN_RES",
	4099:  "WS_CONF_LOCK_REQ",
	4100:  "WS_CONF_LOCK_RES",
	4101:  "WS_CONF_END_REQ",
	4102:  "WS_CONF_END_RES",
	4103:  "WS_CONF_LEAVE_REQ",
	4104:  "WS_CONF_LEAVE_RES",
	4105:  "WS_CONF_RECORD_REQ",
	4106:  "WS_CONF_RECORD_RES",
	4107:  "WS_CONF_EXPEL_REQ",
	4108:  "WS_CONF_EXPEL_RES",
	4109:  "WS_CONF_RENAME_REQ",
	4111:  "WS_CONF_ASSIGN_HOST_REQ",
	4113:  "WS_CONF_PUT_ON_HOLD_REQ",
	4115:  "WS_CONF_SET_MUTE_UPON_ENTRY_REQ",
	4117:  "WS_CONF_SET_HOLD_UPON_ENTRY_REQ",
	4119:  "WS_CONF_INVITE_CRC_DEVICE_REQ",
	4120:  "WS_CONF_INVITE_CRC_DEVICE_RES",
	4121:  "WS_CONF_CANCEL_INVITE_CRC_DEVICE_REQ",
	4122:  "WS_CONF_CANCEL_INVITE_CRC_DEVICE_RES",
	4123:  "WS_CONF_SET_BROADCAST_REQ",
	4124:  "WS_CONF_SET_BROADCAST_RES",
	4125:  "WS_CONF_CLOSED_CAPTION_REQ",
	4126:  "WS_CONF_CLOSED_CAPTION_RES",
	4127:  "WS_CONF_ALLOW_VIEW_PARTICIPANT_REQ",
	4129:  "WS_CONF_LOWER_ALL_HAND_REQ",
	4131:  "WS_CONF_RAISE_LOWER_HAND_REQ",
	4133:  "WS_CONF_RECLAIM_HOST_REQ",
	4135:  "WS_CONF_CHAT_REQ",
	4137:  "WS_CONF_ASSIGN_CC_REQ",
	4141:  "WS_CONF_CHAT_PRIVILEDGE_REQ",
	4143:  "WS_CONF_FEEDBACK_REQ",
	4145:  "WS_CONF_FEEDBACK_CLEAR_REQ",
	4147:  "WS_CONF_ALLOW_UNMUTE_VIDEO_REQ",
	4149:  "WS_CONF_ALLOW_UNMUTE_AUDIO_REQ",
	4151:  "WS_CONF_ALLOW_RAISE_HAND_REQ",
	4153:  "WS_CONF_PANELIST_VOTE_REQ",
	4173:  "WS_CONF_BO_TOKEN_REQ",
	4174:  "WS_CONF_BO_TOKEN_RES",
	4175:  "WS_CONF_BO_START_REQ",
	4177:  "WS_CONF_BO_STOP_REQ",
	4179:  "WS_CONF_BO_ASSIGN_REQ",
	4181:  "WS_CONF_BO_SWITCH_REQ",
	4183:  "WS_CONF_BO_WANT_JOIN_REQ",
	4185:  "WS_CONF_BO_LEAVE_REQ",
	4187:  "WS_CONF_BO_BROADCAST_REQ",
	4189:  "WS_CONF_BO_HELP_REQ",
	4191:  "WS_CONF_BO_HELP_RESULT_REQ",
	4193:  "WS_CONF_BO_JOIN_REQ",
	4194:  "WS_CONF_BO_JOIN_RES",
	4163:  "WS_CONF_ALLOW_PARTICIPANT_RENAME_REQ",
	4171:  "WS_CONF_ALLOW_MESSAGE_FEEDBACK_NOTIFY_REQ",
	4195:  "WS_CONF_REVOKE_COHOST_REQ",
	4197:  "WS_CONF_PLAY_CHIME_OPEN_CLOSE_REQ",
	4199:  "WS_CONF_ADMIT_ALL_SILENT_USERS_REQ",
	4201:  "WS_CONF_BIND_UNBIND_TELE_USR_REQ",
	4203:  "WS_CONF_ALLOW_QA_AUTO_REPLY_REQ",
	4205:  "WS_CONF_EXPEL_ATTENDEE_REQ",
	4206:  "WS_CONF_EXPEL_ATTENDEE_RES",
	4207:  "WS_CONF_PRACTICE_SESSION_REQ",
	4208:  "WS_CONF_PRACTICE_SESSION_RES",
	4209:  "WS_CONF_ROLE_CHANGE_REQ",
	4210:  "WS_CONF_ROLE_CHANGE_RES",
	4211:  "WS_CONF_BO_TOKEN_BATCH_REQ",
	4213:  "WS_CONF_BO_PRE_ASSIGN_REQ",
	4214:  "WS_CONF_BO_PRE_ASSIGN_RES",
	4217:  "WS_CONF_CHANGE_MULTI_PIN_PRIVILGE_REQ",
	4219:  "WS_CONF_SET_GROUP_LAYOUT",
	4215:  "WS_CONF_HOST_KEY_REQ",
	4216:  "WS_CONF_HOST_KEY_RES",
	4222:  "WS_CONF_AVATAR_PERMISSION_CHANGED",
	4229:  "WS_CONF_SUSPEND_MEETING",
	4230:  "WS_CONF_SUSPEND_MEETING_REQ_RESULT",
	8193:  "WS_AUDIO_MUTE_REQ",
	8194:  "WS_AUDIO_MUTE_RES",
	8195:  "WS_AUDIO_DROP_REQ",
	8196:  "WS_AUDIO_DROP_RES",
	8197:  "WS_AUDIO_DIALOUT_REQ",
	8198:  "WS_AUDIO_DIALOUT_RES",
	8199:  "WS_AUDIO_CANCEL_DIALOUT_REQ",
	8200:  "WS_AUDIO_CANCEL_DIALOUT_RES",
	8201:  "WS_AUDIO_MUTEALL_REQ",
	8202:  "WS_AUDIO_MUTEALL_RES",
	8204:  "WS_AUDIO_ALLOW_TALK_REQ",
	8205:  "WS_AUDIO_ALLOW_TALK_RES",
	7937:  "WS_CONF_ROSTER_INDICATION",
	7938:  "WS_CONF_ATTRIBUTE_INDICATION",
	7939:  "WS_CONF_END_INDICATION",
	7940:  "WS_CONF_HOST_CHANGE_INDICATION",
	7941:  "WS_CONF_COHOST_CHANGE_INDICATION",
	7942:  "WS_CONF_HOLD_CHANGE_INDICATION",
	7943:  "WS_CONF_CLOSED_CAPTION_INDICATION",
	7944:  "WS_CONF_CHAT_INDICATION",
	7945:  "WS_CONF_OPTION_INDICATION",
	7946:  "WS_CONF_KV_UPDATE_INDICATION",
	7947:  "WS_CONF_LOCAL_RECORD_INDICATION",
	8203:  "WS_AUDIO_VOIP_JOIN_CHANNEL_REQ",
	7949:  "WS_CONF_BO_COMMAND_INDICATION",
	7950:  "WS_CONF_BO_ATTRIBUTE_INDICATION",
	7951:  "WS_CONF_ADMIT_ALL_SILENT_USERS_INDICATION",
	7952:  "WS_CONF_BIND_UNBIND_INDICATION",
	7953:  "WS_CONF_UPDATE_MEETING_TOPIC_INDICATION",
	7954:  "WS_CONF_DC_REGION_INDICATION",
	7955:  "WS_CONF_CAN_ADMIT_WHEN_NOHOST_PRESENT_INDICATION",
	7958:  "WS_CONF_GROUP_LAYOUT_INDICATION",
	12033: "WS_AUDIO_ASN_INDICATION",
	12034: "WS_AUDIO_MUTE_INDICATION",
	12035: "WS_AUDIO_SSRC_INDICATION",
	12036: "WS_AUDIO_ALLOW_TALK_INDICATION",
	12037: "WS_AUDIO_SSRC_ASK_UNMUTE_INDICATION",
	12038: "WS_WEBINAR_VIEW_ONLY_TELEPHONY_INDICATION",
	16129: "WS_VIDEO_ACTIVE_INDICATION",
	16131: "WS_VIDEO_SSRC_INDICATION",
	16133: "WS_VIDEO_MUTE_INDICATION",
	16135: "WS_VIDEO_LEADERSHIP_INDICATION",
	12289: "WS_VIDEO_SUBSCRIBE_REQ",
	12291: "WS_VIDEO_UNSUBSCRIBE_REQ",
	12293: "WS_VIDEO_KEY_FRAME_REQ",
	12295: "WS_VIDEO_NETWORK_FEEDBACK",
	12297: "WS_VIDEO_MUTE_VIDEO_REQ",
	12299: "WS_VIDEO_SPOTLIGHT_VIDEO_REQ",
	16385: "WS_SHARING_PAUSE_REQ",
	16387: "WS_SHARING_RESUME_REQ",
	20225: "WS_SHARING_STATUS_INDICATION",
	20226: "WS_SHARING_SIZE_CHANGE_INDICATION",
	4155:  "WS_CONF_ALLOW_ANONYMOUS_QUESTION_REQ",
	4157:  "WS_CONF_ALLOW_VIEW_ALL_QUESTION_REQ",
	4159:  "WS_CONF_ALLOW_UPVOTE_QUESTION_REQ",
	4161:  "WS_CONF_ALLOW_COMMENT_QUESTION_REQ",
	16389: "WS_SHARING_REMOTE_CONTROL_REQ",
	16391: "WS_SHARING_REMOTE_CONTROL_INDICATION",
	16393: "WS_SHARING_REMOTE_CONTROLLER_GRAB",
	16395: "WS_SHARING_REMOTE_CONTROLLER_GRAB_INDICATION",
	16415: "WS_SHARING_SUBSCRIBE_REQ",
	16417: "WS_SHARING_UNSUBSCRIBE_REQ",
	4167:  "WS_MEETING_RWG_CONNECT_TIME",
	24322: "EVT_TYPE_WS_VIDEO_DATACHANNEL_ANSWER",
	24321: "DATA_CHANNEL_SEND_OFFER_TO_RWG",
	4223:  "WS_CONF_FOLLOW_HOST_REQ",
	7957:  "WS_CONF_DRAG_LAYOUT_INDICATION",
	4218:  "WS_CONF_SET_DRAG_LAYOUT",
	4227:  "WS_CONF_LIVE_TRANSCRIPTION_ON_OFF_REQ",
	4228:  "WS_CONF_LIVE_TRANSCRIPTION_ON_OFF_RES",
	7959:  "WS_CONF_LIVE_TRANSCRIPTION_STATUS_INDICATION",
	4165:  "WS_CONF_POLLING_REQ",
	4224:  "WS_CONF_POLLING_USER_ACTION_REQ",
	4225:  "WS_CONF_POLLING_USER_ACTION_ERROR",
	4226:  "WS_CONF_POLLING_SET_POLLING_TOKEN",
	// OTHER
	4169:  "WS_CONF_LOCK_SHARE_REQ",
	16409: "WS_CONF_SET_SHARE_STATUS_REQ",
}
