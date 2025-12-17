package go_core_rest_api

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ReadConfig() {
}
