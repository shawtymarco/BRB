package utils

import (
	"math"
	"strings"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/text/cases"
	language2 "golang.org/x/text/language"
)

func Distance(v1 mgl64.Vec3, v2 mgl64.Vec3) float64 {
	return math.Sqrt(math.Pow(v2.X()-v1.X(), 2) + math.Pow(v2.Y()-v1.Y(), 2) + math.Pow(v2.Z()-v1.Z(), 2))
}

func HorizontalDistance(v1 mgl64.Vec3, v2 mgl64.Vec3) float64 {
	return math.Sqrt(math.Pow(v2.X()-v1.X(), 2) + math.Pow(v2.Z()-v1.Z(), 2))
}

func DirectionFromVector(v mgl64.Vec3) mgl64.Vec3 {
	var x, z float64
	if v.X() > v.Z() {
		if v.X() < 0 {
			x = -1
			z = 0
		} else {
			x = 1
			z = 0
		}
	} else {
		if v.Z() < 0 {
			x = 0
			z = -1
		} else {
			x = 0
			z = 1
		}
	}
	return mgl64.Vec3{x, v.Y(), z}
}

func ItemDisplay(it item.Stack) string {
	if it.CustomName() != "" {
		return it.CustomName()
	}
	name, _ := it.Item().EncodeItem()
	name = strings.Split(name, ":")[1]
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	name = cases.Title(language2.English).String(name)
	return name
}

func VecSetY(v mgl64.Vec3, y float64) mgl64.Vec3 {
	return mgl64.Vec3{v.X(), y, v.Z()}
}

var ZeroBox = cube.Box(-0.001, 0, -0.001, 0.001, 0.001, 0.001)
var BlockBox = cube.Box(-0.5, 0, -0.5, 0.5, 1, 0.5)
