package protocol

import "errors"

var (
	ErrNoData        = errors.New("no data was passed")
	ErrInvalidLength = errors.New("invalid length")
)
