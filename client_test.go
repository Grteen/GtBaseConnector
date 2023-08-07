package GtBase

import (
	"context"
	"testing"

	"github.com/grteen/gtBaseconnector/opt"
	"github.com/grteen/gtBaseconnector/pkg"
	"github.com/grteen/gtBaseconnector/pool"
	"github.com/grteen/gtBaseconnector/utils"
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

func TestCommand(t *testing.T) {
	addr := "127.0.0.1:9877"
	gdb := NewClient(&opt.Option{
		Addr: addr,
	})

	err := gdb.Set("Key", "Val").Err()
	if err != nil {
		t.Errorf(err.Error())
	}

	val := gdb.Get("Key").Result()
	if val != "Val" {
		t.Errorf("should get %v but got %v", "Val", val)
	}

	err = gdb.Del("Key").Err()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = gdb.Get("Key").Err()
	if err != pkg.GtBaseNil {
		t.Errorf("should get %v but got %v", pkg.GtBaseNil, err)
	}
}
