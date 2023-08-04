package pool

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pkg"
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

func TestReadResp(t *testing.T) {
	addr := "127.0.0.1:4567"
	opt := &opt.Option{
		Addr: addr,
	}
	opt.Init()
	pool := NewConnPool(opt, opt.Dialer)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		resp := []byte{4, 0, 0, 0}
		resp = append(resp, []byte("Pang")...)
		resp = append(resp, []byte(pkg.CommandSep)...)
		if err := utils.ListenCheckReqAndWriteBack(cancel, addr, []byte("Ping"), resp); err != nil {
			t.Errorf(err.Error())
		}
	}()

	<-ctx.Done()

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
		result, err := cn.ReadResp()
		if err != nil {
			return err
		}

		if !utils.EqualByteSlice(result, []byte("Pang")) {
			t.Errorf("should get Pang but got %v", result)
		}

		return nil
	}); err != nil {
		t.Errorf(err.Error())
	}
}
