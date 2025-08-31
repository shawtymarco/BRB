package command

import (
	"fmt"
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"strconv"
	"time"

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
	Indicator

	SetRole
	ChangeCape
	ResetStats
	GameMode
	Sudo
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
	database.Owner:     {Sudo, GameMode},
	database.Manager:   {SetRole},
	database.Admin:     {ResetStats, ChangeCape},
	database.Moderator: {Alias, Ban, Indicator},
	database.Helper:    {Mute},
	database.Prime:     {Nick, ClaimELO},
	database.Premium:   {Fly, Spectate},
	database.Media:     {},
	database.Booster:   {},
	database.Player:    {Indicator},
}

type ArgumentPlayer string

func (ArgumentPlayer) Type() string {
	return "player"
}

func (ArgumentPlayer) Options(_ cmd.Source) []string {
	var players []string
	for _, name := range server.Players {
		players = append(players, name)
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

type Duration string

func (Duration) Type() string {
	return "duration"
}

func (Duration) Options(_ cmd.Source) []string {
	durs := []string{"permanent"}
	for i := 1; i <= 180; i++ {
		durs = append(durs, fmt.Sprintf("%vm", i))
		durs = append(durs, fmt.Sprintf("%vh", i))
		durs = append(durs, fmt.Sprintf("%vd", i))
	}
	return durs
}

func (arg Duration) Parse() time.Duration {
	dur := string(arg)
	numPart := dur[:len(dur)-1]
	unit := dur[len(dur)-1]

	value := utils.Panics(strconv.Atoi(numPart))

	switch unit {
	case 'm':
		return time.Duration(value) * time.Minute
	case 'h':
		return time.Duration(value) * time.Hour
	case 'd':
		return time.Duration(value) * 24 * time.Hour
	default:
		return 0
	}
}
