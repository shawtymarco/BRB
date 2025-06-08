package server

import (
	"server/server/database"
	"server/server/worldmanager"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server"
)

var (
	Database     database.Database
	MCServer     *server.Server
	Config       Server
	WorldManager *worldmanager.Manager
)

var (
	BotMark   *player.Player
	BotSam    *player.Player
	BotSteven *player.Player
)
