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
)
