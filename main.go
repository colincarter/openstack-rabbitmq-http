package main

import (
	"log"
	"os"

	"bitbucket.org/colincarter/openstack-rabbitmq-http/commands"

	"gopkg.in/urfave/cli.v1"
)

var cliApp *cli.App

func main() {
	cliApp = cli.NewApp()
	cliApp.Name = "openstack-http-rabbitmq"
	cliApp.Usage = "Openstack to RabbitMQ HTTP Broker"
	cliApp.Author = "Colin Carter"
	cliApp.Email = "ccarter@ukcloud.com"
	cliApp.Version = "0.0.0"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "listen, l",
			Value:  "0.0.0.0",
			Usage:  "listen address",
			EnvVar: "HTTP_LISTEN",
		},
		cli.IntFlag{
			Name:   "port, p",
			Value:  3000,
			Usage:  "port to listen on",
			EnvVar: "HTTP_PORT",
		},
		cli.StringFlag{
			Name:   "rabbit-host, r",
			Value:  "localhost",
			Usage:  "hostname of rabbitmq host",
			EnvVar: "RABBIT_HOST",
		},
		cli.IntFlag{
			Name:   "rabbit-port, o",
			Value:  5672,
			Usage:  "rabbitmq port",
			EnvVar: "RABBIT_PORT",
		},
		cli.StringFlag{
			Name:   "rabbit-exchange, e",
			Value:  "/",
			Usage:  "rabbitmq exchange",
			EnvVar: "RABBIT_EXCHANGE",
		},
		cli.StringFlag{
			Name:   "rabbit-user, u",
			Value:  "guest",
			Usage:  "rabbitmq username",
			EnvVar: "RABBIT_USER",
		},
		cli.StringFlag{
			Name:   "rabbit-password, a",
			Value:  "guest",
			Usage:  "rabbitmq password",
			EnvVar: "RABBIT_PASSWORD",
		},
		cli.IntFlag{
			Name:   "concurrency, c",
			Value:  5,
			Usage:  "number of rabbit producers",
			EnvVar: "CONCURRENCY",
		},
		cli.StringFlag{
			Name:   "failure-dir, f",
			Value:  "",
			Usage:  "place to store failures",
			EnvVar: "FAILURE-DIR",
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run webserver",
			Action: func(c *cli.Context) error {
				return commands.Server(c)
			},
			Flags: flags,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
