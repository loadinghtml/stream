package structs

type Stream struct {
	API        StreamAPI `json:"api"`
	DNS        StreamDNS `json:"dns"`
	TCP        StreamTCP `json:"tcp"`
	Strict     bool      `json:"strict"`
	WhiteList  []string  `json:"whitelist"`
	DomainList []string  `json:"domainlist"`
}

type StreamAPI struct {
	Listen string   `json:"listen"`
	Secret []string `json:"secret"`
}

type StreamDNS struct {
	Strict   bool   `json:"strict"`
	Listen   string `json:"listen"`
	Address  string `json:"address"`
	Upstream string `json:"upstream"`
}

type StreamTCP struct {
	HTP []string `json:"htp"`
	TLS []string `json:"tls"`
}
