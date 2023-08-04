package GtBase

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pool"
)

type Client struct {
	*BaseClient

	pool *pool.ConnPool
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
