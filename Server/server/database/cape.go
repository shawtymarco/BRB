package database

import (
	"image"
	"image/png"
	"os"
	"server/server/utils"
)

var capeRegistry = map[string]Cape{}

func RegisterCape(cape Cape) {
	capeRegistry[cape.Identifier()] = cape
}

func AllCapes() map[string]Cape {
	return capeRegistry
}

type Cape interface {
	Identifier() string
	Name() string
	ImagePath() string
}

func GetCapeByIdentifier(t string) (Cape, bool) {
	cape, ok := capeRegistry[t]
	return cape, ok
}

func CapeAsImage(cape Cape) image.Image {
	return utils.Panics(png.Decode(utils.Panics(os.Open(cape.ImagePath()))))
}

func CapeAsBytes(cape Cape) []uint8 {
	return CapeAsImage(cape).(*image.NRGBA).Pix
}
