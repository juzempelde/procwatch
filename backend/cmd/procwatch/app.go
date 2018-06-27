package main

import (
	"gopkg.in/urfave/cli.v2"
)

func createApp() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "agent",
				Action: func(ctx *cli.Context) error {
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "server-rpc-addr",
					},
				},
			},
			&cli.Command{
				Name: "server",
				Action: func(ctx *cli.Context) error {
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "rpc-addr",
					},
				},
			},
		},
	}
}
