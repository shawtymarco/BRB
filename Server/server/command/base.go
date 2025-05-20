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
	WorldEdit Permission = iota
	GiveRank
	Gamemode
	GameForceStart
	ChangeNickname
	ChatColor
)

func (p Permission) Test(src cmd.Source) bool {
	if pl, ok := src.(*player.Player); ok {
		u := user.LookupPlayer(pl)
		if u != nil {
			pRank := u.Data.GroupSettings.Rank()
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
					return text.Colourf(language.Translate(pl).Global.Commands.Error.Permission, r.Icon(), r.Prefix())
				}
			}
		}
	}
	return text.Colourf("<red>Something went wrong...</red>")
}

var rankPermissions = map[database.Rank][]Permission{
	database.Owner:     {},
	database.Council:   {WorldEdit, GiveRank, Gamemode, GameForceStart},
	database.Sheriff:   {},
	database.Mason:     {},
	database.Developer: {},
	database.Artist:    {},
	database.Partner:   {},
	database.Platinum:  {},
	database.Gold:      {},
	database.Silver:    {ChangeNickname},
	database.VIP:       {ChatColor},
	database.Player:    {},
}
