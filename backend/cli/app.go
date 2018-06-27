package cli

import (
	"gopkg.in/urfave/cli.v2"
)

// CreateApp creates an application which can be fed with os.Args.
func CreateApp(commands Commands) *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "agent",
				Action: func(ctx *cli.Context) error {
					agent := commands.Agent()
					agent.SetServerRPCAddr(ctx.String("server-rpc-addr"))
					return agent.Exec()
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
					server := commands.Server()
					server.SetAddr(ctx.String(serverRPCAddr))
					return server.Exec()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: serverRPCAddr,
					},
				},
			},
		},
	}
}

const (
	agentServerRPCAddr = "server-rpc-addr"
	serverRPCAddr      = "rpc-addr"
)

// Execer can be executed. If that fails, an error is returned.
type Execer interface {
	Exec() error
}

// Agent executes as an agent.
type Agent interface {
	Execer

	SetServerRPCAddr(addr string)
}

// Server executes as a server.
type Server interface {
	Execer

	SetAddr(addr string)
}

// Commands provides all the commands the application can run.
type Commands interface {
	Agent() Agent
	Server() Server
}
