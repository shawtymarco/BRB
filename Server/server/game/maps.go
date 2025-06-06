package game

import (
	"github.com/go-gl/mathgl/mgl64"
)

var Maps map[string]MapData

type MapData struct {
	Name            string
	Type            TypeGame
	Author          string
	Mode            int
	SpawnPoint      mgl64.Vec3
	TeamSpawnPoints []mgl64.Vec3
}
