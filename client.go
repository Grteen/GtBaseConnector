package GtBase

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pool"
	"GtBase-Connector/utils"
)

type BaseClient struct {
	opt *opt.Option

	pool *pool.ConnPool
}

func (c *BaseClient) WithConn(fn func(cn *pool.GtBaseConn) error) error {
	cn, err := c.pool.GetConn()
	if err != nil {
		return err
	}

	defer c.pool.PushIdle(cn)

	err = fn(cn)
	return err
}

func WriteReq(cn *pool.GtBaseConn, cmder *Cmder) error {
	bts := utils.EncodeFieldsToGtBasePacket(cmder.fields)

	return cn.Write(bts)
}

func (c *BaseClient) process(cmder *Cmder) (*Cmder, error) {
	if err := c.WithConn(func(cn *pool.GtBaseConn) error {
		if err := cn.WithWriteFunc(c.opt.WriteTimeOut, func() error {
			return WriteReq(cn, cmder)
		}); err != nil {
			return err
		}

		if err := cn.WithReadFunc(c.opt.ReadTimeOut, func() error {
			bts, err := cn.ReadResp()
			if err != nil {
				return err
			}
			cmder.reply = bts

			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		cmder.setFnErr(err)
		return cmder, err
	}

	return cmder, nil
}

type Client struct {
	*BaseClient
	CmdAble
}

func NewClient(opt *opt.Option) *Client {
	opt.Init()

	c := &Client{
		BaseClient: &BaseClient{
			opt: opt,
		},
	}
	c.init()
	c.pool = pool.NewConnPool(opt, c.opt.Dialer)

	return c
}

func (c *Client) init() {
	c.CmdAble = c.process
}
