package pool

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/grteen/gtbaseconnector/opt"
	"github.com/grteen/gtbaseconnector/pkg"
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

func (p *ConnPool) removeConn(cn *GtBaseConn) {
	for i, c := range p.conns {
		if cn == c {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled {
				p.poolSize--
			}
			break
		}
	}
}

func (p *ConnPool) removeConnLock(cn *GtBaseConn) {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	p.removeConn(cn)
}

func (p *ConnPool) closeConn(cn *GtBaseConn) error {
	return cn.Close()
}

// remove Conn from p.conns and close it
func (p *ConnPool) freeConn(cn *GtBaseConn) error {
	p.removeConnLock(cn)
	return p.closeConn(cn)
}

// PushIdle will freeConn if cn is bad
func (p *ConnPool) PushIdle(cn *GtBaseConn) error {
	if p.closed() {
		return pkg.ClosedError
	}

	if cn.IsBad() {
		return p.freeConn(cn)
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	p.idleConns = append(p.idleConns, cn)
	p.idleLen++

	return nil
}

func (p *ConnPool) PopIdle() (*GtBaseConn, error) {
	if p.closed() {
		return nil, pkg.ClosedError
	}

	n := len(p.idleConns)
	if n == 0 {
		return nil, nil
	}

	cn := p.idleConns[0]
	copy(p.idleConns, p.idleConns[1:])
	p.idleConns = p.idleConns[:n-1]
	p.idleLen--

	return cn, nil
}

func (p *ConnPool) GetConn() (*GtBaseConn, error) {
	if p.closed() {
		return nil, pkg.ClosedError
	}

	for {
		p.connsMu.Lock()
		cn, err := p.PopIdle()
		p.connsMu.Unlock()
		if err != nil {
			return nil, err
		}

		// no connection so create one
		if cn == nil {
			break
		}

		// ToDo check connection health and continue
		return cn, nil
	}

	cn, err := p.newConn(true)
	if err != nil {
		cn.Close()
		return nil, err
	}

	return cn, nil
}
