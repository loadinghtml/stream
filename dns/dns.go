package dns

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aiocloud/stream/api"

	mdns "github.com/miekg/dns"
)

var (
	StrictDNS bool = false

	mux       *mdns.ServeMux
	tcpSocket *mdns.Server
	udpSocket *mdns.Server
	dialer    = net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: time.Second * 10,

		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", api.Get().DNS.Upstream)
			},
		},
	}
)

func Dial(network, address string) (net.Conn, error) {
	return dialer.Dial(network, address)
}

func Listen(addr string) {
	if addr == "" {
		return
	}

	list := api.Get().DomainList
	mux = mdns.NewServeMux()
	mux.HandleFunc("in-addr.arpa.", handleServerName)
	for i := 0; i < len(list); i++ {
		mux.HandleFunc(mdns.Fqdn(list[i]), handleDomain)
	}
	mux.HandleFunc(".", handleOther)

	tcpSocket = &mdns.Server{Net: "tcp", Addr: addr, Handler: mux}
	udpSocket = &mdns.Server{Net: "udp", Addr: addr, Handler: mux}

	go func() { log.Fatalf("[Stream][DNS][TCP] %v", tcpSocket.ListenAndServe()) }()
	go func() { log.Fatalf("[Stream][DNS][UDP] %v", udpSocket.ListenAndServe()) }()

	log.Println("[Stream][DNS] Started")
}

func handleServerName(w mdns.ResponseWriter, r *mdns.Msg) {
	if StrictDNS {
		if checked, err := api.CheckIP(w.RemoteAddr()); err != nil {
			log.Printf("[Stream][DNS][api.CheckIP] %v", err)
			return
		} else if !checked {
			return
		}
	}

	m := new(mdns.Msg)
	m.SetReply(r)

	for i := 0; i < len(r.Question); i++ {
		rr, err := mdns.NewRR(fmt.Sprintf("%s PTR aioCloud", r.Question[i].Name))
		if err != nil {
			log.Println(err)
			return
		}

		m.Answer = append(m.Answer, rr)
	}

	_ = w.WriteMsg(m)
}

func handleDomain(w mdns.ResponseWriter, r *mdns.Msg) {
	if StrictDNS {
		if checked, err := api.CheckIP(w.RemoteAddr()); err != nil {
			log.Printf("[Stream][DNS][api.CheckIP] %v", err)
			return
		} else if !checked {
			return
		}
	}

	m := new(mdns.Msg)
	m.SetReply(r)

	for i := 0; i < len(r.Question); i++ {
		rr, err := mdns.NewRR(fmt.Sprintf("%s A %s", r.Question[i].Name, api.Get().DNS.Address))
		if err != nil {
			log.Println(err)
			return
		}

		m.Answer = append(m.Answer, rr)
	}

	_ = w.WriteMsg(m)
}

func handleOther(w mdns.ResponseWriter, r *mdns.Msg) {
	if StrictDNS {
		if checked, err := api.CheckIP(w.RemoteAddr()); err != nil {
			log.Printf("[Stream][DNS][api.CheckIP] %v", err)
			return
		} else if !checked {
			return
		}
	}

	m, err := mdns.Exchange(r, api.Get().DNS.Upstream)
	if err != nil {
		return
	}

	_ = w.WriteMsg(m)
}
