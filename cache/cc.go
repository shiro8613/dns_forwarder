package cache

import (
	"dns_forwarder/config"
	"fmt"
	"github.com/miekg/dns"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

type Cache struct {
	m  map[string]Value
	mu sync.RWMutex
}

type Value struct {
	Server dns.RR
	Time   int64
}

var c Cache

func Init() {
	c = Cache{
		m: make(map[string]Value),
	}
}

func SetAutoRemove() error {
	t := config.GetConfig().CacheTime * 10
	cr := cron.New()
	_, err := cr.AddFunc(fmt.Sprintf("@every %ds", t), autoRemoveFunc)
	if err != nil {
		return err
	}

	cr.Start()

	return nil
}

func autoRemoveFunc() {
	c.mu.Lock()
	for k, m := range c.m {
		if time.Now().After(time.UnixMilli(m.Time)) {
			delete(c.m, k)
		}
	}
	c.mu.Unlock()
}

func Is(name string) bool {
	c.mu.RLock()
	s, ok := c.m[name]
	c.mu.RUnlock()
	if ok {
		return time.UnixMilli(s.Time).After(time.Now())
	}

	return false
}

func Get(name string) dns.RR {
	c.mu.RLock()
	s := c.m[name]
	c.mu.RUnlock()
	return s.Server
}

func Set(name string, server dns.RR) {
	limit := config.GetConfig().CacheTime
	limitTime := time.Now().Add(time.Duration(limit) * time.Second)

	c.mu.Lock()
	c.m[name] = Value{
		Server: server,
		Time:   limitTime.UnixMilli(),
	}
	c.mu.Unlock()
}
