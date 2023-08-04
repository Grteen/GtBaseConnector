package utils

import (
	"context"
	"fmt"
	"net"
)

func ListenCheckReqAndWriteBack(cancel context.CancelFunc, addr string, req []byte, resp []byte) error {
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Listen Success
	cancel()

	cn, err := listner.Accept()
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	n, err := cn.Read(buf)
	if err != nil {
		return err
	}

	if !EqualByteSlice(buf[:n], req) {
		return fmt.Errorf("should get %v but got %v", req, buf[:n])
	}

	_, err = cn.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
