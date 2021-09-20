package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/aiocloud/stream/api"
	"github.com/aiocloud/stream/dns"
	"github.com/aiocloud/stream/mitm"
)

var (
	Path string
)

func main() {
	flag.StringVar(&Path, "c", "/etc/stream.json", "Path")
	flag.Parse()

	if err := api.Load(Path); err != nil {
		log.Fatalf("[Stream] api.Load: %v", err)
	}

	if err := api.UpdateIPv4(); err != nil {
		log.Printf("[Stream] api.UpdateIPv4: %v", err)
	}

	if err := api.UpdateIPv6(); err != nil {
		log.Printf("[Stream] api.UpdateIPv6: %v", err)
	}

	if api.CurrentIPv4 == "" && api.CurrentIPv6 == "" {
		log.Fatalln("[Stream] Get current ip address failed")
	}

	if err := api.UpdateRule(); err != nil {
		log.Fatalf("[Stream] Update rule failed: %v", err)
	}

	api.Run()
	dns.Run()

	for _, i := range api.StreamData.TCP.TLS {
		go mitm.ListenTLS(i)
	}

	for _, i := range api.StreamData.TCP.HTTP {
		go mitm.ListenHTTP(i)
	}

	go UpdateIP()
	go UpdateRule()

	log.Printf("[Stream] IPv4: %s IPv6: %s", api.CurrentIPv4, api.CurrentIPv6)
	log.Println("[Stream] Started")

	for {
		time.Sleep(time.Minute * 10)

		runtime.GC()

		stats := new(runtime.MemStats)
		runtime.ReadMemStats(stats)
		log.Printf("[Stream][GC] CPU Fraction %f", stats.GCCPUFraction)
		log.Printf("[Stream][GC] Obtained %dMB", stats.Sys/1024/1024)
		log.Printf("[Stream][GC] Assigned %dMB", stats.Alloc/1024/1024)
		log.Printf("[Stream][GC] Routine %d", runtime.NumGoroutine())
	}
}

func UpdateIP() {
	for {
		time.Sleep(time.Second * 120)

		if err := api.UpdateIPv4(); err != nil {
			log.Printf("[Stream] api.UpdateIPv4: %v", err)
		}

		if err := api.UpdateIPv6(); err != nil {
			log.Printf("[Stream] api.UpdateIPv6: %v", err)
		}

		log.Printf("[Stream] IPv4: %s IPv6: %s", api.CurrentIPv4, api.CurrentIPv6)
	}
}

func UpdateRule() {
	for {
		time.Sleep(time.Second * 86400)

		if err := api.UpdateRule(); err != nil {
			log.Printf("[Stream] Update rule failed: %v", err)
		}
	}
}
