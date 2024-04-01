package config

import "regexp"

func (s Server) GetV4() []string {
	return *s.ToV4
}

func (s Server) GetV6() []string {
	return *s.ToV6
}

func splitV4V6(c *Config) error {
	if err := split(c.Forwards, &c.ForwardsV4, &c.ForwardsV6); err != nil {
		return err
	}

	c.Servers = mapSplit(c.Servers)

	return nil
}

func mapSplit(m *map[string]Server) *map[string]Server {
	m1 := make(map[string]Server)

	for k, s := range *m {
		v4 := new([]string)
		v6 := new([]string)

		if err := split(s.To, v4, v6); err != nil {
			return nil
		}

		s.ToV4 = v4
		s.ToV6 = v6

		m1[k] = s
	}

	return &m1
}

func split(l []string, v4 *[]string, v6 *[]string) error {
	v4exp, err := regexp.Compile(ipv4Regex)
	if err != nil {
		return err
	}

	for _, s := range l {
		if v4exp.MatchString(s) {
			*v4 = append(*v4, s)
		} else {
			*v6 = append(*v6, s)
		}
	}
	return nil
}
