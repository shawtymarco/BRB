package server

import (
	"path"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var IP = "BRBW.NET"

type Server struct {
	Prefix    string
	Languages map[string][]string
	Pvp       struct {
		Force           float64
		Height          float64
		HitRegistration int
	}
	Hub struct {
		SpawnPoint mgl64.Vec3
	}
}

func DefaultConfig() server.UserConfig {
	c := server.DefaultConfig()
	c.Network.Address = ":19135"

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
