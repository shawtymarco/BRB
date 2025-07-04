package intercept

import "github.com/df-mc/dragonfly/server"

var srv *server.Server

func Start(s *server.Server) {
	srv = s
}
