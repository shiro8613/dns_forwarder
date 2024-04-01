package main

import (
	"dns_forwarder/cache"
	"dns_forwarder/config"
	"dns_forwarder/logger"
	"dns_forwarder/proxy"
	"fmt"
	"github.com/miekg/dns"
	"sync"
)

func init() {
	logger.SetupLogger("./log")
	if err := config.Load(); err != nil {
		logger.Panic(err)
	}
	cache.Init()
}

func main() {
	conf := config.GetConfig()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		logger.Info(fmt.Sprintf("(v4)Listen on %s", conf.Bind4))
		logger.Panic(dns.ListenAndServe(conf.Bind4, "udp", proxy.Handler("v4")))
		wg.Done()
	}()

	if conf.Bind6 != "" {
		wg.Add(1)
		go func() {
			logger.Info(fmt.Sprintf("(v6)Listen on %s", conf.Bind6))
			logger.Panic(dns.ListenAndServe(conf.Bind6, "udp", proxy.Handler("v6")))
			wg.Done()
		}()
	}

	logger.Info("Server is Started")
	wg.Wait()
}
