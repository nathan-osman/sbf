package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/sbf/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sbf",
		Usage: "Microservice for sending large files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "address to listen on",
				Value:   ":http",
				EnvVars: []string{"ADDR"},
			},
			&cli.StringFlag{
				Name:    "dir",
				Usage:   "directory for storing files",
				Value:   ".",
				EnvVars: []string{"DIR"},
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "username for the admin",
				EnvVars: []string{"USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "password for the admin",
				EnvVars: []string{"PASSWORD"},
			},
		},
		Action: func(c *cli.Context) error {

			// Create the tmp directory if it doesn't exist
			if err := os.MkdirAll(os.TempDir(), 1777); err != nil {
				return err
			}

			// Create & initialize the server
			s, err := server.New(&server.Config{
				Addr:     c.String("addr"),
				Dir:      c.String("dir"),
				Username: c.String("username"),
				Password: c.String("password"),
			})
			if err != nil {
				return err
			}
			defer s.Close()

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
