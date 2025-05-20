package maps

import (
	"github.com/go-gl/mathgl/mgl64"
)

var MapPath string
var MapConfig MapCollection

type MapCollection struct {
	TowerWars map[string]TWMapConf
}

type TWMapConf struct {
	Name         string
	Author       string
	Mode         int
	WLSpawnPoint mgl64.Vec3
	Direction    string

	WhiteSpawn         mgl64.Vec3
	WhiteSummon        mgl64.Vec3
	WhiteCastle        mgl64.Vec3
	WhiteLeftTower     mgl64.Vec3
	WhiteRightTower    mgl64.Vec3
	WhiteLeftPath      []mgl64.Vec3
	WhiteRightPath     []mgl64.Vec3
	WhiteShopMagicBook mgl64.Vec3

	WhiteShopIceShield   mgl64.Vec3
	WhiteShopFireShield  mgl64.Vec3
	WhiteShopEarthShield mgl64.Vec3
	WhiteShopMagicShield mgl64.Vec3

	BlackSpawn         mgl64.Vec3
	BlackSummon        mgl64.Vec3
	BlackCastle        mgl64.Vec3
	BlackLeftTower     mgl64.Vec3
	BlackRightTower    mgl64.Vec3
	BlackLeftPath      []mgl64.Vec3
	BlackRightPath     []mgl64.Vec3
	BlackShopMagicBook mgl64.Vec3

	BlackShopIceShield   mgl64.Vec3
	BlackShopFireShield  mgl64.Vec3
	BlackShopEarthShield mgl64.Vec3
	BlackShopMagicShield mgl64.Vec3
}
