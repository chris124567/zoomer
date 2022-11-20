package zoom

import (
	"errors"
)

var (
	ErrParticipantExists  = errors.New("participant already exists")
	ErrParticipantMissing = errors.New("participant is missing")
	ErrSsrcMissing        = errors.New("ssrc is missing")
)

type Participant struct {
	userId      int
	ssrcs       []int
	secretNonce []byte
}

type ZoomParticipantRoster struct {
	sharedMeetingKey []byte
	participants     map[ /*userId*/ int]*Participant
}

func NewParticipantRoster() *ZoomParticipantRoster {
	return &ZoomParticipantRoster{
		participants: make(map[int]*Participant),
	}
}
func (roster *ZoomParticipantRoster) AddParticipant(userId int, zoomId string) error {
	secretNonce, err := ZoomEscapedBase64Decode(zoomId)
	if err != nil {
		return err
	}

	participant := roster.participants[userId]
	if participant != nil {
		return ErrParticipantExists
	}

	newParticipant := &Participant{
		userId:      userId,
		ssrcs:       make([]int, 0),
		secretNonce: secretNonce,
	}
	roster.participants[userId] = newParticipant

	return nil
}

func (roster *ZoomParticipantRoster) AddSsrcForParticipant(userId int, ssrc int) error {
	participant := roster.participants[userId]
	if participant == nil {
		return ErrParticipantMissing
	}

	participant.ssrcs = append(participant.ssrcs, ssrc)
	return nil
}

func (roster *ZoomParticipantRoster) GetSecretNonceForSSRC(ssrcNeedle int) ([]byte, error) {
	var secretNonce []byte
	for _, participant := range roster.participants {
		found := false
		for _, ssrcHay := range participant.ssrcs {
			if ssrcNeedle == ssrcHay {
				found = true
				secretNonce = participant.secretNonce
				break
			}
		}
		if found {
			break
		}
	}

	if len(secretNonce) == 0 {
		return nil, ErrSsrcMissing
	}

	return secretNonce, nil
}

func (roster *ZoomParticipantRoster) SetSharedMeetingKey(sharedMeetingKey []byte) {
	roster.sharedMeetingKey = sharedMeetingKey
}

func (roster *ZoomParticipantRoster) GetSharedMeetingKey() []byte {
	return roster.sharedMeetingKey
}
