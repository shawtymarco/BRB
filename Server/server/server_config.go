package server

import (
	"path"

	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Server struct {
	Prefix    string
	WTPrefix  string
	Languages map[string][]string
	Pvp       struct {
		VerticalKB      float64
		HorizontalKB    float64
		HitRegistration int
	}
}

func DefaultConfig() server.UserConfig {
	c := server.DefaultConfig()
	c.Network.Address = ":19132"

	c.Server.Name = text.Colourf("<dark-red>BRB</dark-red>")
	c.Server.DisableJoinQuitMessages = true
	c.Server.AuthEnabled = false
	c.Players.MaxCount = 200
	c.Players.SaveData = false
	c.World.SaveData = true
	c.World.Folder = path.Join(".", "server", "worlds", "hub")
	c.Resources.AutoBuildPack = false
	c.Resources.Required = true
	return c
}
