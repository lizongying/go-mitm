package proxy

import (
	"net"
)

type CustomAddr struct {
	network string
	address string
}

func (a *CustomAddr) Network() string {
	return a.network
}

func (a *CustomAddr) String() string {
	return a.address
}

type Conn struct {
	net.Conn
	addr net.Addr
}

func NewConn(conn net.Conn, addr string) net.Conn {
	return &Conn{
		Conn: conn,
		addr: &CustomAddr{
			address: addr,
		},
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.addr
}
