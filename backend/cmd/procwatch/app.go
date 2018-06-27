package main

import (
	"github.com/juzempelde/procwatch/backend/agent"
	"github.com/juzempelde/procwatch/backend/cli"
	"github.com/juzempelde/procwatch/backend/cli/wrap"
	"github.com/juzempelde/procwatch/backend/server"
)

func createApp() Runner {
	return cli.CreateApp(
		cli.NewCommands(
			func() cli.Agent {
				return wrap.Agent(&agent.Agent{})
			}, func() cli.Server {
				return wrap.Server(&server.Server{})
			},
		),
	)
}
