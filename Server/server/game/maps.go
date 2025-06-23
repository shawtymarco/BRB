package game

import (
	"github.com/go-gl/mathgl/mgl64"
)

var Maps map[string]MapData

type MapData struct {
	Name   string
	Type   TypeGame
	Author string
	Mode   int

	Void int

	SpawnPoint      mgl64.Vec3
	TeamSpawnPoints []mgl64.Vec3

	BedPositions              []mgl64.Vec3
	ShopVillagerPositions     []mgl64.Vec3
	UpgradesVillagerPositions []mgl64.Vec3
	IronGenerators            []mgl64.Vec3
	DiamondGenerators         []mgl64.Vec3
	EmeraldGenerators         []mgl64.Vec3
}
