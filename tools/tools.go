package tools

import (
	"io"
	"net"
	"time"
)

func AddrToIP(addr string) net.IP {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	return net.ParseIP(host)
}

func IPToIPNet(addr net.IP) (n *net.IPNet) {
	if addr.To4() != nil {
		_, n, _ = net.ParseCIDR(addr.String() + "/32")
	} else {
		_, n, _ = net.ParseCIDR(addr.String() + "/128")
	}

	return n
}

func CopyBuffer(client, remote net.Conn) {
	go func() {
		_, _ = io.CopyBuffer(remote, client, make([]byte, 1446))
		_ = client.SetDeadline(time.Now())
		_ = remote.SetDeadline(time.Now())
	}()

	_, _ = io.CopyBuffer(client, remote, make([]byte, 1446))
	_ = client.SetDeadline(time.Now())
	_ = remote.SetDeadline(time.Now())
}
