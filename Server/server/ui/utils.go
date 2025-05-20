package ui

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

const (
	GameSelector = "game.selector"
)

var values = make(map[*world.EntityHandle]map[string]interface{})

func Store[T any](pl *world.EntityHandle, key string, value T) {
	if values[pl] == nil {
		values[pl] = make(map[string]interface{})
	}
	values[pl][key] = value
}

func Load[T any](pl *player.Player, key string) T {
	return values[pl.H()][key].(T)
}

func AddButtonWithValue(pl *player.Player, text, image string, value interface{}) form.Button {
	Store(pl.H(), text, value)
	return form.NewButton(text, image)
}
