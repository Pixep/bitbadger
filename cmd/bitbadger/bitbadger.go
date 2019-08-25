package main

import (
	bitbadger "github.com/Pixep/bitbadger/internal/bitbadger"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	// VERSION stores the current version as string
	VERSION = "v0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "BitBadger"
	app.Version = VERSION
	app.Usage = "A badge generator for BitBucket"
	app.Action = start
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug mode",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "Set the port that the server listens on",
			Value: 34000,
		},
	}

	app.Run(os.Args)
}

func start(c *cli.Context) error {
	if c.NArg() < 2 {
		log.Fatal("Please provide a Username and Password.")
	}

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	config := bitbadger.Config{
		Username: c.Args().Get(0),
		Password: c.Args().Get(1),
	}
	bitbadger.SetConfig(config)

	log.Info("Serving badges as '", config.Username, "'")

	bitbadger.Serve(c.Int("port"))
	return nil
}
