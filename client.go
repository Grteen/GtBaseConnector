package GtBase

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pool"
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

type Client struct {
	*BaseClient
}

func NewClient(opt *opt.Option) *Client {
	opt.Init()

	c := &Client{
		BaseClient: &BaseClient{
			opt: opt,
		},
	}
	c.pool = pool.NewConnPool(opt, c.opt.Dialer)

	return c
}
