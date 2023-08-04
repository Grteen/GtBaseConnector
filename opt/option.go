package opt

import (
	"GtBase-Connector/pkg"
	"net"
	"time"
)

// Option Keeps all settings to set up GtBase connections
type Option struct {
	Addr string

	// Dialer Creates new Tcp Connection
	Dialer func(addr string) (net.Conn, error)

	// 0 is default 3 * time.Second
	// -1 is block
	// -2 is unblock
	DialTimeOut time.Duration
	// 0 is default 3 * time.Second
	// -1 is block
	// -2 is unblock
	ReadTimeOut time.Duration
	// 0 is default 3 * time.Second
	// -1 is block
	// -2 is unblock
	WriteTimeOut time.Duration

	MaxPoolSize int
}

func (opt *Option) Init() {
	if opt.Addr == "" {
		opt.Addr = pkg.DefaultConnectAddress
	}
	if opt.DialTimeOut == 0 {
		opt.DialTimeOut = pkg.DefaultConnectTimeOut
	}
	if opt.Dialer == nil {
		opt.Dialer = NewDialer(opt)
	}
	if opt.MaxPoolSize == 0 {
		opt.MaxPoolSize = pkg.DefaultMaxPoolSize
	}
}

func NewDialer(opt *Option) func(addr string) (net.Conn, error) {
	return func(addr string) (net.Conn, error) {
		netDialer := net.Dialer{
			Timeout: opt.DialTimeOut,
		}

		return netDialer.Dial("tcp", addr)
	}
}
