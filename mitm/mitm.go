package mitm

import (
	"log"
	"net"

	"github.com/aiocloud/stream/api"
)

func ListenHTTP(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[Stream][HTTP][net.Listen] %v", err)
	}
	defer ln.Close()

	log.Printf("[Stream][HTTP][%s] Started", addr)

	for {
		client, err := ln.Accept()
		if err != nil {
			log.Fatalf("[Stream][HTTP][ln.Accept] %v", err)
		}

		if checked, err := api.CheckIP(client.RemoteAddr()); err != nil {
			log.Printf("[Stream][HTTP][api.CheckIP] %v", err)

			_ = client.Close()
			continue
		} else if !checked {
			log.Printf("[Stream][HTTP][api.CheckIP] Ban %s", client.RemoteAddr())

			_ = client.Close()
			continue
		}

		go handleHTTP(client)
	}
}

func ListenTLS(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[Stream][TLS][net.Listen] %v", err)
	}
	defer ln.Close()

	log.Printf("[Stream][TLS][%s] Started", addr)

	for {
		client, err := ln.Accept()
		if err != nil {
			log.Fatalf("[Stream][TLS][ln.Accept] %v", err)
		}

		if checked, err := api.CheckIP(client.RemoteAddr()); err != nil {
			log.Printf("[Stream][TLS][api.CheckIP] %v", err)

			_ = client.Close()
			continue
		} else if !checked {
			log.Printf("[Stream][TLS][api.CheckIP] Ban %s", client.RemoteAddr())

			_ = client.Close()
			continue
		}

		go handleTLS(client)
	}
}
