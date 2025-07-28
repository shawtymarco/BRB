package game

import (
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/item"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Team struct {
	id               int
	originalHandles  []*world.EntityHandle
	activeHandles    []*world.EntityHandle
	spectatorHandles []*world.EntityHandle
	color            string

	Status   Status
	Upgrades struct {
		Sharpness     int
		Protection    int
		Haste         int
		HealPool      int
		GeneratorTier int

		Traps          [3]Trap
		ActiveTrap     Trap
		ActivatedSince time.Time
	}
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

func (t *Team) Colour() string {
	return t.color
}

func (t *Team) WoolColour() item.Colour {
	switch t.color {
	case "red":
		return item.ColourRed()
	case "blue":
		return item.ColourBlue()
	case "green":
		return item.ColourGreen()
	case "yellow":
		return item.ColourYellow()
	}

	return item.ColourBlack()
}

func (t *Team) TrapsCount() (count int) {
	for _, trap := range t.Upgrades.Traps {
		if trap != None {
			count++
		}
	}
	return count
}

func (t *Team) IsTrapsFull() bool {
	return t.TrapsCount() == 3
}

func (t *Team) AddTrap(trap Trap) {
	for i := range t.Upgrades.Traps {
		if t.Upgrades.Traps[i] == None {
			t.Upgrades.Traps[i] = trap
			return
		}
	}

	panic("AddTrap called when trap queue is full")
}

func (t *Team) RemoveTrap() Trap {
	removed := t.Upgrades.Traps[0]
	copy(t.Upgrades.Traps[0:], t.Upgrades.Traps[1:])
	t.Upgrades.Traps[2] = None
	return removed
}

func ColorToID(color string, typeGame TypeGame) int {
	switch typeGame {
	case TypeBedWars:
		switch color {
		case "red":
			return 0
		case "blue":
			return 1
		case "green":
			return 2
		case "yellow":
			return 3
		}
	case TypeBedFight, TypeBuildFFA:
		switch color {
		case "red":
			return 0
		case "blue":
			return 2
		case "green":
			return 1
		case "yellow":
			return 3
		}
	}
	panic("no id for this color")
}

type Status int

const (
	BedExists Status = iota
	BedBroken
	TeamDead
)

type Trap int

func (t Trap) Slot() int {
	switch t {
	case Regular:
		return 14
	case CounterOffensive:
		return 15
	case Alarm:
		return 16
	case MinerFatigue:
		return 23
	default:
		return -1
	}
}

const (
	None Trap = iota
	Regular
	CounterOffensive
	Alarm
	MinerFatigue
)
