package main

import (
	"github.com/juzempelde/procwatch/backend"
	"github.com/juzempelde/procwatch/backend/agent"
	"github.com/juzempelde/procwatch/backend/cli"
	"github.com/juzempelde/procwatch/backend/cli/wrap"
	"github.com/juzempelde/procwatch/backend/ps"
	"github.com/juzempelde/procwatch/backend/server"

	"os"
)

func createApp() Runner {
	return cli.CreateApp(
		cli.NewCommands(
			func() cli.Agent {
				return wrap.Agent(
					&agent.Agent{
						ProcessList:    procwatch.ProcessListFunc(ps.CurrentProcesses),
						HostIDProvider: procwatch.HostIDProviderFunc(os.Hostname),
					},
				)
			}, func() cli.Server {
				return wrap.Server(
					&server.Server{
						Server: &procwatch.Server{
							Devices: &procwatch.Devices{},
						},
					},
				)
			},
		),
	)
}
