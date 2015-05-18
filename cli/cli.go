package cli

import (
	"io"

	"github.com/codegangsta/cli"
)

// Run starts the CLI.
func Run(version string, writer io.Writer, args []string) error {
	app := cli.NewApp()
	app.Name = "Restdis"
	app.Version = version
	app.Usage = "HTTP interface for Redis"
	app.Writer = writer
	app.Action = Start

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file",
			EnvVar: "RESTDIS_CONFIG",
		},
		cli.StringFlag{
			Name:   "bind, b",
			Usage:  "Address to bind the web service on",
			EnvVar: "RESTDIS_BIND",
		},
		cli.StringFlag{
			Name:   "redis, r",
			Usage:  "Redis DSN in the form <host:port>",
			EnvVar: "RESTDIS_REDIS",
		},
	}

	return app.Run(args)
}
