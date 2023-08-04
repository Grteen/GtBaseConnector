package pkg

import "time"

const (
	DefaultConnectAddress string        = "127.0.0.1:9877"
	DefaultConnectTimeOut time.Duration = 10 * time.Second

	CommandSep             string = "\r\n"
	GtBasePacketLengthSize int32  = 4
)
