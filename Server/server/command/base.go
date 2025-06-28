package command

import (
	"fmt"
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/world"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Permission int

const (
	Fly Permission = iota
	Spectate
	Nick
	ClaimELO

	Mute
	Ban
	Alias

	SetRole
	ChangeCape
	ResetStats
)

func (p Permission) Test(src cmd.Source) bool {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u != nil {
			pRank := u.Data.Rank()
			for pRank <= database.Player {
				for _, perm := range rankPermissions[pRank] {
					if perm == p {
						return true
					}
				}
				pRank++
			}
			return false
		}
	}
	return true
}

func (p Permission) PermissionMessage(src cmd.Source) string {
	if pl, ok := src.(*player.Player); ok && !p.Test(src) {
		for r, perms := range rankPermissions {
			for _, perm := range perms {
				if perm == p {
					return text.Colourf(language.Translate(pl).Commands.Error.Permission, r.Prefix())
				}
			}
		}
	}
	return text.Colourf("<red>Something went wrong...</red>")
}

var rankPermissions = map[database.Rank][]Permission{
	database.Owner:     {},
	database.Manager:   {},
	database.Admin:     {SetRole, ChangeCape},
	database.Moderator: {Alias, Ban},
	database.Helper:    {Mute},
	database.Prime:     {Nick, ClaimELO},
	database.Premium:   {Fly, Spectate},
	database.Media:     {},
	database.Booster:   {},
	database.Player:    {ResetStats}, // TODO: PUT BACK TO OWNER PERMS
}

type ArgumentPlayer string

func (ArgumentPlayer) Type() string {
	return "player"
}

func (ArgumentPlayer) Options(src cmd.Source) []string {
	var players []string
	for pl := range server.MCServer.Players(src.(*player.Player).Tx()) {
		players = append(players, pl.Name())
	}
	return players
}

func (arg ArgumentPlayer) ResolveIn(tx *world.Tx) (*player.Player, bool) {
	handle, ok := server.MCServer.PlayerByName(string(arg))
	if !ok {
		return nil, false
	}
	entity, ok := handle.Entity(tx)
	if !ok {
		return nil, false
	}
	pl, ok := entity.(*player.Player)
	return pl, ok
}

func (arg ArgumentPlayer) ExecWithPlayerSafe(currentTx *world.Tx, fn func(tx *world.Tx, pl *player.Player)) error {
	if pl, ok := arg.ResolveIn(currentTx); ok {
		fn(currentTx, pl)
		return nil
	}

	handle, ok := server.MCServer.PlayerByName(string(arg))
	if !ok {
		return fmt.Errorf("player %s not found", string(arg))
	}

	var ran bool
	handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
		p, ok := e.(*player.Player)
		if !ok {
			return
		}
		fn(tx, p)
		ran = true
	})
	if !ran {
		return fmt.Errorf("player %s could not be resolved", string(arg))
	}
	return nil
}
