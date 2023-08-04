package pool

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/utils"
	"context"
	"testing"
)

func TestDial(t *testing.T) {
	addrf := "127.0.0.1:1234"
	addrt := "127.0.0.1:9877"
	opt := &opt.Option{
		Addr: addrf,
	}
	opt.Init()
	pool := NewConnPool(opt, opt.Dialer)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := utils.ListenCheckReqAndWriteBack(cancel, addrt, []byte("Ping"), []byte("Pang")); err != nil {
			t.Errorf(err.Error())
		}
	}()

	<-ctx.Done()

	_, err := pool.newConn(true)
	// should fail
	if err == nil {
		t.Errorf("should get Connection Refused but got %v", err)
	}

	opt.Addr = addrt
	cn, err := pool.newConn(true)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := cn.WithWriteFunc(opt.WriteTimeOut, func() error {
		_, err := cn.conn.Write([]byte("Ping"))
		return err
	}); err != nil {
		t.Errorf(err.Error())
	}

	if err := cn.WithReadFunc(opt.ReadTimeOut, func() error {
		buf := make([]byte, 1024)
		n, err := cn.conn.Read(buf)
		if err != nil {
			return err
		}

		if !utils.EqualByteSlice(buf[:n], []byte("Pang")) {
			t.Errorf("should get Pang but got %v", buf[:n])
		}

		return nil
	}); err != nil {
		t.Errorf(err.Error())
	}
}
