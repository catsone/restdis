package cli

import (
	"os"

	"github.com/catsone/restdis/service"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/zenazn/goji/bind"
)

// Start starts a worker.
func Start(c *cli.Context) {
	config, err := loadConfig(c)
	if err != nil {
		log.Fatalf("Configuration error: %s", err)
	}

	service := service.RestdisService{
		Version: c.App.Version,
		Config:  config,
	}

	err = service.Run()
	if err != nil {
		log.Fatalf("Error starting service: %s", err)
	}
}

// loadConfig loads the service configuration from a config file if one has been specified, or if "config.toml" exists
// in the current directory.
func loadConfig(c *cli.Context) (*service.Config, error) {
	// Config file
	var config service.Config

	file := c.String("config")

	if file == "" {
		if _, err := os.Stat("config.toml"); err == nil {
			file = "config.toml"
		}
	}

	if file != "" {
		_, err := toml.DecodeFile(file, &config)
		if err != nil {
			return nil, err
		}
	}

	// Bind
	bindOpt := c.String("bind")

	if bindOpt == "" {
		bindOpt = bind.Sniff()
	}

	if bindOpt == "" {
		bindOpt = ":7631"
	}

	config.Bind = bindOpt

	// Redis
	redis := c.String("redis")

	if redis != "" {
		config.RedisAddr = redis
	}

	if config.RedisAddr == "" {
		config.RedisAddr = "127.0.0.1:6379"
	}

	return &config, nil
}
