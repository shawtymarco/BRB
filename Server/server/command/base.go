package command

import (
	"server/server/database"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Permission int

const (
	GiveRank Permission = iota
	GameMode
	Warp
	Join
)

func (p Permission) Test(src cmd.Source) bool {
	if pl, ok := src.(*player.Player); ok {
		u := user.LookupPlayer(pl)
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
	database.Owner:        {GameMode, GiveRank},
	database.Manager:      {},
	database.Admin:        {},
	database.Moderator:    {},
	database.Helper:       {},
	database.Prime:        {},
	database.Premium:      {},
	database.MediaPartner: {},
	database.Booster:      {},
	database.Player:       {Join, Warp},
}
