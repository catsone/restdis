package service

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const contextErrorStatusCode = "errorCode"

type header struct {
	Name  string
	Value string
}

// Config is the configuration for the Restdis web service.
type Config struct {
	Bind      string `toml:"bind"`
	RedisAddr string `toml:"redis"`
	Headers   []header
}

// RestdisService represents the Restdis web service.
type RestdisService struct {
	Version string
	Config  *Config
}

// Run starts the Restdis web service.
func (s *RestdisService) Run() error {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", s.Config.RedisAddr)

			// TODO: Redis auth.
			return c, err
		},
	}

	resource := &RestdisResource{version: s.Version, pool: pool}

	r := web.New()

	r.Use(middleware.EnvInit)
	r.Use(middleware.Logger)
	r.Use(ErrorHandler)
	r.Use(Headers(s.Config))
	r.Use(middleware.AutomaticOptions)

	r.Get("/", resource.Default)
	r.Post("/", resource.Default)
	r.Get("/*", resource.RedisCommand)
	r.Post("/*", resource.RedisCommand)

	listener := bind.Socket(s.Config.Bind)
	log.Infof("Starting Restdis on %s", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()

	graceful.PreHook(func() {
		log.Info("Restdis received signal, gracefully stopping")
	})
	graceful.PostHook(func() {
		log.Info("Restdis stopped")
	})

	return graceful.Serve(listener, r)
}
