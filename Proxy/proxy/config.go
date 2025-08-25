package proxy

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

type Config struct {
	Connection struct {
		LocalAddress  string
		RemoteAddress string
	}
}

func ReadConfig() Config {
	c := Config{}
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("create config: %v", err)
		}
		data, err := toml.Marshal(c)
		if err != nil {
			log.Fatalf("encode default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("write default config: %v", err)
		}
		_ = f.Close()
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("decode config: %v", err)
	}
	if c.Connection.LocalAddress == "" {
		c.Connection.LocalAddress = "0.0.0.0:19132"
	}
	data, _ = toml.Marshal(c)
	if err := os.WriteFile("config.toml", data, 0644); err != nil {
		log.Fatalf("write config: %v", err)
	}
	return c
}
