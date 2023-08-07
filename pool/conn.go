package pool

import (
	"net"
	"time"

	"github.com/grteen/gtbaseconnector/pkg"
	"github.com/grteen/gtbaseconnector/utils"
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

func (conn *GtBaseConn) Write(bts []byte) error {
	_, err := conn.conn.Write(bts)
	return err
}

func (conn *GtBaseConn) ReadResp() ([]byte, error) {
	respLen, err := conn.readRespLen()
	if err != nil {
		return nil, err
	}

	return conn.readResp(respLen)
}

func (conn *GtBaseConn) readRespLen() (int32, error) {
	buf := make([]byte, 4)
	n, err := conn.conn.Read(buf)
	if err != nil {
		return -1, err
	}

	return utils.EncodeBytesSmallEndToint32(buf[:n]), nil
}

func (conn *GtBaseConn) readResp(respLen int32) ([]byte, error) {
	result := make([]byte, 0, respLen)
	restn := respLen

	for restn > 0 {
		buf := make([]byte, respLen)
		n, err := conn.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		restn -= int32(n)
		result = append(result, buf[:n]...)
	}

	// Read the \r\n
	buf := make([]byte, 2)
	n, err := conn.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if !utils.EqualByteSlice(buf[:n], []byte(pkg.CommandSep)) {
		return nil, pkg.InvalidRespError
	}

	return result, nil
}

func (conn *GtBaseConn) IsBad() bool {
	return !conn.pooled
}
