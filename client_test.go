package GtBase

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pkg"
	"GtBase-Connector/pool"
	"GtBase-Connector/utils"
	"context"
	"testing"
)

func TestWithConn(t *testing.T) {
	addr := "127.0.0.1:4561"
	c := NewClient(&opt.Option{
		Addr: addr,
	})

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

	if err := c.WithConn(func(cn *pool.GtBaseConn) error {
		if err := cn.WithWriteFunc(c.opt.WriteTimeOut, func() error {
			err := cn.Write([]byte("Ping"))
			return err
		}); err != nil {
			return err
		}

		if err := cn.WithReadFunc(c.opt.ReadTimeOut, func() error {
			bts, err := cn.ReadResp()
			if err != nil {
				return err
			}

			if !utils.EqualByteSlice(bts, []byte("Pang")) {
				t.Errorf("should get %v but got %v", []byte("Pang"), bts)
			}

			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Errorf(err.Error())
	}
}
