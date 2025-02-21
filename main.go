package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sbf",
		Usage: "Microservice for sending large files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "admin-password",
				Usage:   "password for the admin",
				Value:   ":http",
				EnvVars: []string{"ADMIN_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "admin-username",
				Usage:   "username for the admin",
				Value:   ":http",
				EnvVars: []string{"ADMIN_USER"},
			},
			&cli.StringFlag{
				Name:    "server-addr",
				Usage:   "address to listen on",
				Value:   ":http",
				EnvVars: []string{"SERVER_ADDR"},
			},
		},
		Action: func(c *cli.Context) error {

			// Wait for SIGINT or SIGTERM
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
	}
}
