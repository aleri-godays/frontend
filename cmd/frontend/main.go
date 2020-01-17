package main

import (
	"github.com/aleri-godays/frontend/internal"
	"github.com/aleri-godays/frontend/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	version     = "dev-snapshot"
	serviceName = "frontend"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = serviceName
	cliApp.Version = version

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "log level: TRACE, DEBUG, INFO, WARN, ERROR",
			EnvVar: "LOG_LEVEL",
			Value:  "DEBUG",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "port to server",
			EnvVar: "PORT",
			Value:  5000,
		},
		cli.StringFlag{
			Name:   "project-uri",
			Usage:  "uri for project service",
			EnvVar: "PROJECT_URI",
			Value:  "http://localhost:5010",
		},
		cli.StringFlag{
			Name:   "timetracking-uri",
			Usage:  "uri for timetracking service",
			EnvVar: "TIMETRACKING_URI",
			Value:  "http://localhost:5020",
		},
		cli.StringFlag{
			Name:   "jwt-secret",
			Usage:  "jwt-secret",
			EnvVar: "JWT_SECRET",
			Value:  "aslkdhasljkdhasjdh",
		},
		cli.StringFlag{
			Name:   "oauth-client-id",
			Usage:  "oauth-client-id",
			EnvVar: "OAUTH_CLIENT_ID",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "oauth-client-secret",
			Usage:  "oauth-client-secret",
			EnvVar: "OAUTH_CLIENT_SECRET",
			Value:  "",
		},
	}
	cliApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run",
			Action: func(c *cli.Context) error {
				conf := getConfig(c)
				app := internal.NewApp(conf)
				app.Run()
				return nil
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize app")
	}
}

func getConfig(c *cli.Context) *config.Config {
	conf := &config.Config{
		ServiceName:       serviceName,
		Version:           version,
		LogLevel:          c.GlobalString("log-level"),
		HTTPPort:          c.GlobalInt("port"),
		ProjectURI:        c.GlobalString("project-uri"),
		TimeTrackingURI:   c.GlobalString("timetracking-uri"),
		OAuthClientID:     c.GlobalString("oauth-client-id"),
		OAuthClientSecret: c.GlobalString("oauth-client-secret"),
		JWTSecret:         c.GlobalString("jwt-secret"),
	}

	return conf
}
