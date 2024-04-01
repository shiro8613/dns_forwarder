package config

type Config struct {
	Bind4      string   `yaml:"bind4"`
	Bind6      string   `yaml:"bind6"`
	CacheTime  int      `yaml:"cache_time"`
	Forwards   []string `yaml:"forwards"`
	ForwardsV4 []string
	ForwardsV6 []string
	Servers    *map[string]Server `yaml:"servers"`
}

type Server struct {
	To       []string `yaml:"to"`
	ToV4     *[]string
	ToV6     *[]string
	Priority int `yaml:"priority"`
}
