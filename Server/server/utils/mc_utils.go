package utils

import (
	"math"
	"strings"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
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

func ColorDisplay(color string) string {
	return text.Colourf("<%v>%v</%v>", color, cases.Title(language2.English, cases.NoLower).String(strings.Replace(color, "-", " ", -1)), color)
}

func ColorSlice(slice []string) []string {
	var coloredSlice []string
	for _, s := range slice {
		coloredSlice = append(coloredSlice, text.Colourf(s))
	}
	return coloredSlice
}

func VecNoY(v mgl64.Vec3) mgl64.Vec3 {
	return mgl64.Vec3{v.X(), 0, v.Z()}
}

func RegisterItem(i world.Item) {
	world.RegisterItem(i)

	group := "Custom Blocks"

	if customItem, ok := i.(interface {
		Group() string
	}); ok {
		group = customItem.Group()
	}

	creative.RegisterItem(creative.Item{Stack: item.NewStack(i, 1), Group: group})
}

var ZeroBox = cube.Box(-0.001, 0, -0.001, 0.001, 0.001, 0.001)
var BlockBox = cube.Box(-0.5, 0, -0.5, 0.5, 1, 0.5)
