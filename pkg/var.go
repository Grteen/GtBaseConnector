package pkg

import (
	"errors"
	"time"
)

var (
	NoDeadLine = time.Time{}

	// all operations on the closed client will return this error
	ClosedError = errors.New("client is closed")
	// resp from GtBase is invalid will return this error
	InvalidRespError = errors.New("invalid response")

	NilReply []byte = []byte{0, 0, 0, 78, 105, 108}
	// resp from GtBase is nil will return this error
	GtBaseNil = errors.New("nil")
)
