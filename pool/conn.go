package pool

import (
	"GtBase-Connector/pkg"
	"net"
	"time"
)

type GtBaseConn struct {
	conn net.Conn

	// figure out if this Connection is Pooled or not
	pooled bool
}

func NewGtBaseConn(conn net.Conn) *GtBaseConn {
	gtcon := &GtBaseConn{
		conn: conn,
	}

	return gtcon
}

func (conn *GtBaseConn) calDeadLine(timeout time.Duration) time.Time {
	tm := time.Now()

	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	if timeout > 0 {
		return tm
	}

	return pkg.NoDeadLine
}

func (conn *GtBaseConn) WithWriteFunc(timeout time.Duration, fn func() error) error {
	if timeout > 0 {
		if err := conn.conn.SetWriteDeadline(conn.calDeadLine(timeout)); err != nil {
			return err
		}
	}

	return fn()
}

func (conn *GtBaseConn) WithReadFunc(timeout time.Duration, fn func() error) error {
	if timeout > 0 {
		if err := conn.conn.SetReadDeadline(conn.calDeadLine(timeout)); err != nil {
			return err
		}
	}

	return fn()
}

func (conn *GtBaseConn) Close() error {
	return conn.conn.Close()
}
