package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/catsone/restdis/cli"
)

// Current Restdis version.
const Version = "0.0.1"

func main() {
	err := cli.Run(Version, os.Stdout, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
