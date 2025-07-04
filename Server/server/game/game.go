package game

import (
	"slices"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

type Game struct {
	id         uuid.UUID
	world      *world.World
	stage      Stage
	teams      []*Team
	teamColor  string
	spectators []*world.EntityHandle
}

func NewGame(id uuid.UUID, world *world.World, teamColor string) *Game {
	return &Game{id: id, world: world, teamColor: teamColor}
}

func (g *Game) ID() uuid.UUID {
	return g.id
}

func (g *Game) World() *world.World {
	return g.world
}

func (g *Game) Stage() Stage {
	return g.stage
}

func (g *Game) SetStage(stage Stage) {
	g.stage = stage
}

func (g *Game) Teams() []*Team {
	return g.teams
}

func (g *Game) OriginalPlayers() (handles []*world.EntityHandle) {
	for _, team := range g.teams {
		handles = append(handles, team.originalHandles...)
	}
	return handles
}

func (g *Game) ActivePlayers() (handles []*world.EntityHandle) {
	for _, team := range g.teams {
		handles = append(handles, team.activeHandles...)
	}
	return handles
}

func (g *Game) ForEachOriginalPlayer(f func(pl *player.Player), tx *world.Tx) {
	for _, e := range g.OriginalPlayers() {
		if e, ok := e.Entity(tx); ok {
			if pl, ok := e.(*player.Player); ok {
				f(pl)
			}
		}
	}
}

func (g *Game) ForEachActivePlayer(f func(pl *player.Player), tx *world.Tx) {
	for _, e := range g.ActivePlayers() {
		if e, ok := e.Entity(tx); ok {
			if pl, ok := e.(*player.Player); ok {
				f(pl)
			}
		}
	}
}

func (g *Game) PlayerTeam(pl *player.Player) *Team {
	for _, team := range g.teams {
		if team.Contains(pl) {
			return team
		}
	}
	return nil
}

func (g *Game) AddPlayerToTeam(pl *player.Player, teamSize int) {
	for _, team := range g.teams {
		if len(team.originalHandles) < teamSize {
			team.AddPlayer(pl)
			return
		}
	}

	var newTeam *Team
	if g.teamColor == "" {
		teamColors := []string{"red", "blue", "green", "yellow"}
		newTeam = &Team{id: len(g.teams), color: teamColors[len(g.teams)]}
	} else {
		newTeam = &Team{color: g.teamColor}
	}
	newTeam.AddPlayer(pl)
	g.teams = append(g.teams, newTeam)
}

func (g *Game) RemovePlayerFromTeam(pl *player.Player) {
	t := g.PlayerTeam(pl)
	t.RemovePlayerFromOriginal(pl)
	t.RemovePlayerFromActive(pl)
}

func (g *Game) IsSpectator(pl *player.Player) bool {
	return slices.Contains(g.spectators, pl.H())
}

func (g *Game) AddSpectator(pl *player.Player) {
	g.spectators = append(g.spectators, pl.H())
}

func (g *Game) RemoveSpectator(pl *player.Player) {
	for i, h := range g.spectators {
		if h == pl.H() {
			g.spectators = append(g.spectators[:i], g.spectators[i+1:]...)
			break
		}
	}
}

func (g *Game) EnemyWith(p1 *player.Player, p2 *player.Player) bool {
	return g.PlayerTeam(p1) != g.PlayerTeam(p2)
}

func (g *Game) WinningTeam() *Team {
	var res *Team
	for _, team := range g.teams {
		if len(team.activeHandles) > 0 {
			if res == nil {
				res = team
			} else {
				return nil
			}
		}
	}

	return res
}

type Essentials interface {
	Type() TypeGame
	Maps() []string
	MapConfig() MapData
	Handler() player.Handler
	Reward(player *player.Player, tx *world.Tx) (before, after int, mvp bool)
	Punish(player *player.Player, tx *world.Tx) (before, after int)
}

const (
	TypeBedWars  TypeGame = "bedwars"
	TypeBedFight TypeGame = "bedfight"
	TypeBuildFFA TypeGame = "buildffa"
)

type TypeGame string

func (t TypeGame) Title() string {
	switch t {
	case TypeBedWars:
		return "BedWars"
	case TypeBedFight:
		return "BedFight"
	case TypeBuildFFA:
		return "BuildFFA"
	}

	return ""
}
