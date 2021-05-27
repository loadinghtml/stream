package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/aiocloud/stream/structs"
	"github.com/aiocloud/stream/tools"
	"github.com/gin-gonic/gin"
	"github.com/yl2chen/cidranger"
)

var (
	info  structs.Stream
	cidr  cidranger.Ranger
	mutex sync.RWMutex
)

func Get() *structs.Stream {
	return &info
}

func Load(name string) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &info); err != nil {
		return err
	}

	cidr = cidranger.NewPCTrieRanger()
	for i := 0; i < len(info.WhiteList); i++ {
		_, n, err := net.ParseCIDR(info.WhiteList[i])
		if err != nil {
			return err
		}

		cidr.Insert(cidranger.NewBasicRangerEntry(*n))
	}

	return nil
}

func Boot() {
	if info.API.Listen == "" {
		channel := make(chan os.Signal, 1)
		signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
		<-channel
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/aio", func(c *gin.Context) {
		a := tools.AddrToIP(c.Request.RemoteAddr)

		mutex.RLock()
		if checked, _ := cidr.Contains(a); checked {
			mutex.RUnlock()

			c.String(http.StatusOK, fmt.Sprintf("DONE %s", a.String()))
			return
		}
		mutex.RUnlock()

		s := c.Query("secret")
		n := tools.IPToIPNet(a)
		for i := 0; i < len(info.API.Secret); i++ {
			if s == info.API.Secret[i] {
				mutex.Lock()
				cidr.Insert(cidranger.NewBasicRangerEntry(*n))
				mutex.Unlock()

				c.String(http.StatusOK, fmt.Sprintf("DONE %s", a.String()))
				return
			}
		}

		c.String(http.StatusForbidden, "FAIL")
	})
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	log.Fatalf("[Stream][API][r.Run] %v", r.Run(info.API.Listen))
}

func CheckIP(name net.Addr) (bool, error) {
	host, _, err := net.SplitHostPort(name.String())
	if err != nil {
		return false, err
	}

	mutex.RLock()
	defer mutex.RUnlock()

	return cidr.Contains(net.ParseIP(host))
}

func CheckDoamin(name string) bool {
	if !info.Strict {
		return true
	}

	for i := 0; i < len(info.DomainList); i++ {
		if strings.HasSuffix(name, info.DomainList[i]) {
			return true
		}
	}

	return false
}
