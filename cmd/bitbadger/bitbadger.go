package main

import (
	"errors"
	"os"
	"time"

	bitbadger "github.com/Pixep/bitbadger/internal/bitbadger"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// VERSION stores the current version as string
	VERSION = "v0.1.2"
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
		cli.BoolFlag{
			Name:  "insecure, i",
			Usage: "Enable insecure HTTP, without TLS",
		},
		cli.StringFlag{
			Name:  "cert, c",
			Usage: "Path to TLS certificate",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "Path to TLS private key",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "Set the port that the server listens on",
			Value: 34000,
		},
		cli.IntFlag{
			Name:  "cachevalidity",
			Usage: "Set for how long the requests should be cached in minutes",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "maxcached",
			Usage: "Set the maximum number of cached requests",
			Value: 100,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	os.Exit(0)
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

	bitbadger.SetCachePolicy(bitbadger.CachePolicy{
		ValidityDuration: time.Duration(c.Int("cachevalidity")) * time.Minute,
		MaxCachedResults: c.Int("maxcached"),
	})

	log.Info("Serving badges as '", config.Username, "'")

	if c.Bool("insecure") {
		log.Info("Running in HTTP-mode")
		return bitbadger.ServeWithHTTP(c.Int("port"))
	}

	certFile := c.String("cert")
	if certFile == "" {
		log.Error("No certificate provided.")
		return errors.New("No certificate was provided")
	}

	keyFile := c.String("key")
	if keyFile == "" {
		return errors.New("No private key was provided")
	}

	return bitbadger.ServeWithHTTPS(c.Int("port"), certFile, keyFile)
}
