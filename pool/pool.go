package pool

import (
	"GtBase-Connector/opt"
	"GtBase-Connector/pkg"
	"net"
	"sync"
	"sync/atomic"
)

type PoolConfig struct {
	dial func() (net.Conn, error)

	maxPoolSize int
}

type ConnPool struct {
	cfg *PoolConfig

	connsMu   sync.Mutex
	conns     []*GtBaseConn
	idleConns []*GtBaseConn

	poolSize int
	idleLen  int

	_closed uint32 // atomic
}

func NewConnPool(opt *opt.Option, dialer func(string) (net.Conn, error)) *ConnPool {
	p := &ConnPool{
		cfg: &PoolConfig{
			dial: func() (net.Conn, error) {
				return dialer(opt.Addr)
			},
			maxPoolSize: opt.MaxPoolSize,
		},

		conns:     make([]*GtBaseConn, 0, opt.MaxPoolSize),
		idleConns: make([]*GtBaseConn, 0, opt.MaxPoolSize),
	}

	return p
}

func (p *ConnPool) newConn(wantPool bool) (*GtBaseConn, error) {
	cn, err := p.dialConn()
	if err != nil {
		return nil, err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	if p.closed() {
		cn.Close()
		return nil, pkg.ClosedError
	}

	p.conns = append(p.conns, cn)
	if wantPool {
		if p.poolSize >= p.cfg.maxPoolSize {
			cn.pooled = false
		} else {
			cn.pooled = true
			p.poolSize++
		}
	}

	return cn, nil
}

func (p *ConnPool) dialConn() (*GtBaseConn, error) {
	if p.closed() {
		return nil, pkg.ClosedError
	}

	netConn, err := p.cfg.dial()
	if err != nil {
		return nil, err
	}

	conn := NewGtBaseConn(netConn)
	return conn, nil
}

func (p *ConnPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}
