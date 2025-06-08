package capes

import "path"

type CreeperCape struct {
}

func (c CreeperCape) Identifier() string {
	return "creeper"
}

func (c CreeperCape) Name() string {
	return "Creeper Cape"
}

func (c CreeperCape) ImagePath() string {
	return path.Join(".", "config", "capes", "creeper_cape.png")
}
