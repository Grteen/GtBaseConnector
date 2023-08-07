package GtBase

import (
	"github.com/grteen/gtbaseconnector/pkg"
	"github.com/grteen/gtbaseconnector/utils"
)

type Cmder struct {
	fields [][]byte
	reply  []byte
	fnErr  error
}

func (c *Cmder) Clear() {
	c.fields = make([][]byte, 0)
	c.reply = make([]byte, 0)
}

func (c *Cmder) setFnErr(err error) {
	c.fnErr = err
}

func (c *Cmder) Err() error {
	return c.fnErr
}

func (c *Cmder) Result() string {
	if c.reply == nil {
		return ""
	}
	return string(c.reply)
}

func (c *Cmder) CheckReply() {
	if utils.EqualByteSlice(c.reply, pkg.NilReply) {
		c.fnErr = pkg.GtBaseNil
	}
}

type CmdAble func(*Cmder) (*Cmder, error)

func (c CmdAble) set(key string, value string, cmder *Cmder) *Cmder {
	cmder.Clear()
	cmder.fields = append(cmder.fields, []byte(pkg.CommandSet))
	cmder.fields = append(cmder.fields, []byte(key))
	cmder.fields = append(cmder.fields, []byte(value))

	status, _ := c(cmder)
	return status
}

func (c CmdAble) Set(key string, value string) *Cmder {
	cmder := &Cmder{}
	return c.set(key, value, cmder)
}

func (c CmdAble) get(key string, cmder *Cmder) *Cmder {
	cmder.Clear()
	cmder.fields = append(cmder.fields, []byte(pkg.CommandGet))
	cmder.fields = append(cmder.fields, []byte(key))

	status, _ := c(cmder)
	return status
}

func (c CmdAble) Get(key string) *Cmder {
	cmder := &Cmder{}
	return c.get(key, cmder)
}

func (c CmdAble) del(key string, cmder *Cmder) *Cmder {
	cmder.Clear()
	cmder.fields = append(cmder.fields, []byte(pkg.CommandDel))
	cmder.fields = append(cmder.fields, []byte(key))

	status, _ := c(cmder)
	return status
}

func (c CmdAble) Del(key string) *Cmder {
	cmder := &Cmder{}
	return c.del(key, cmder)
}
