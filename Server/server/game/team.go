package game

import (
	"slices"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Team struct {
	id               int
	originalHandles  []*world.EntityHandle
	activeHandles    []*world.EntityHandle
	spectatorHandles []*world.EntityHandle
	color            string

	Status Status
}

func (t *Team) ID() int {
	return t.id
}

func (t *Team) CountActivePlayers() int {
	return len(t.activeHandles)
}

func (t *Team) ForEachPlayer(tx *world.Tx, f func(pl *player.Player)) {
	for _, e := range t.originalHandles {
		if e, ok := e.Entity(tx); ok {
			if pl, ok := e.(*player.Player); ok {
				f(pl)
			}
		}
	}
}

func (t *Team) Contains(pl *player.Player) bool {
	return slices.Contains(t.originalHandles, pl.H())
}

func (t *Team) AddPlayer(pl *player.Player) {
	t.originalHandles = append(t.originalHandles, pl.H())
	t.activeHandles = append(t.activeHandles, pl.H())
}

func (t *Team) RemovePlayerFromOriginal(pl *player.Player) {
	for i, h := range t.originalHandles {
		if h == pl.H() {
			t.originalHandles = append(t.originalHandles[:i], t.originalHandles[i+1:]...)
			break
		}
	}
}

func (t *Team) RemovePlayerFromActive(pl *player.Player) {
	for i, h := range t.activeHandles {
		if h == pl.H() {
			t.activeHandles = append(t.activeHandles[:i], t.activeHandles[i+1:]...)
			break
		}
	}
}

func (t *Team) RemovePlayerFromSpectate(pl *player.Player) {
	for i, h := range t.spectatorHandles {
		if h == pl.H() {
			t.spectatorHandles = append(t.spectatorHandles[:i], t.spectatorHandles[i+1:]...)
			break
		}
	}
}

func (t *Team) MovePlayerToSpectate(pl *player.Player) {
	t.RemovePlayerFromActive(pl)
	t.spectatorHandles = append(t.spectatorHandles, pl.H())
}

func (t *Team) MovePlayerFromSpectate(pl *player.Player) {
	t.RemovePlayerFromSpectate(pl)
	t.activeHandles = append(t.activeHandles, pl.H())
}

func (t *Team) Color() string {
	return t.color
}

type Status int

const (
	BedExists Status = iota
	BedBroken
	TeamDead
)
