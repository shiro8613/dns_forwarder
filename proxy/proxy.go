package proxy

import (
	"dns_forwarder/cache"
	"dns_forwarder/config"
	"dns_forwarder/logger"
	"dns_forwarder/utils"
	"github.com/miekg/dns"
	"strings"
)

func Handler(t string) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		conf := config.GetConfig()
		for _, q := range m.Question {
			name := q.Name
			if cache.Is(name) {
				s := cache.Get(name)
				ms := new(dns.Msg)
				ms.SetReply(m)
				ms.Answer = append(ms.Answer, s)
				if err := w.WriteMsg(ms); err != nil {
					logger.Error(err)
				}
				return
			}

			server := findServer(*conf.Servers, name)
			if server == nil {
				if err := forwarder(w, m, t); err != nil {
					logger.Error(err)
				}
				return
			}

			serv := server.GetV4()
			if t == "v6" && len(server.GetV6()) > 0 {
				serv = server.GetV6()
			}

			if err := forward(w, m, serv); err != nil {
				logger.Error(err)
			}
		}
	}
}

func findServer(m map[string]config.Server, name string) *config.Server {
	servers := make([]config.Server, 0)
	var server *config.Server = nil

	for k, v := range m {
		if strings.HasSuffix(name, k+".") {
			servers = append(servers, v)
		}
	}

	if len(servers) < 0 {
		return nil
	}

	if len(servers) == 1 {
		return &servers[0]
	}

	for _, s := range servers {
		if server == nil {
			server = &s
			continue
		}
		if server.Priority > s.Priority {
			server = &s
			continue
		}
	}

	return server
}

func forwarder(w dns.ResponseWriter, m *dns.Msg, t string) error {
	conf := config.GetConfig()
	serv := conf.ForwardsV4
	if t == "v6 " && len(conf.ForwardsV6) > 0 {
		serv = conf.ForwardsV6
	}

	return forward(w, m, serv)
}

func forward(w dns.ResponseWriter, m *dns.Msg, fList []string) error {
	client := new(dns.Client)

	for i, s := range fList {
		if !strings.Contains(s, ":") {
			s += ":53"
		}

		r, _, err := client.Exchange(m, s)
		if len(fList) == i && err != nil {
			if err := w.WriteMsg(utils.NoRecordAnswer(m)); err != nil {
				return err
			}
		}
		if err != nil {
			continue
		}

		if err := w.WriteMsg(r); err != nil {
			return err
		} else {
			for ii, ss := range r.Answer {
				if len(r.Question) > ii {
					cache.Set(r.Question[ii].Name, ss)
				}
			}
			break
		}
	}

	return nil
}
