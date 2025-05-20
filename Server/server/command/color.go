package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ColorCommand struct{}

func (p ColorCommand) Run(_ cmd.Source, o *cmd.Output, _ *world.Tx) {
	colors := []string{"black", "dark-blue", "dark-green", "dark-aqua", "dark-red", "dark-purple", "gold", "grey", "dark-grey", "blue", "green", "aqua", "red", "purple", "yellow", "white", "dark-yellow", "quartz", "iron", "netherite", "redstone", "copper", "gold", "emerald", "diamond", "lapis", "amethyst", "obfuscated", "bold", "italic"}
	for _, color := range colors {
		o.Printf("%v: %v", color, text.Colourf("<%v>test TEST</%v>", color, color))
	}
}
