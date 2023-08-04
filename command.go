package GtBase

import "GtBase-Connector/pkg"

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
