package config

import (
	"flag"
	"os"
)

type Config struct {
	URLServer string
}

func New() *Config {
	c := Config{}
	c.parseFlags()
	c.parseEnv()
	return &c
}

func (c *Config) parseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "s", "localhost:8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")
	flag.Parse()
	c.URLServer = URLServer
}

func (c *Config) parseEnv() {
	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}
}
