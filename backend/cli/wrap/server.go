package wrap

import (
	"github.com/juzempelde/procwatch/backend/cli"
	"github.com/juzempelde/procwatch/backend/server"
)

// Server wraps a server into something the CLI package can use.
func Server(srv *server.Server) cli.Server {
	return &cliServer{
		server: srv,
	}
}

type cliServer struct {
	server *server.Server
}

func (srv *cliServer) Exec() error {
	return srv.server.Run()
}

func (srv *cliServer) SetAddr(addr string) {
	srv.server.RPCAddr = addr
}
