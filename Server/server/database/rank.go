package database

import "github.com/sandertv/gophertunnel/minecraft/text"

type Rank int

const (
	Owner Rank = iota
	Developer
	Manager
	Admin
	Designer
	Moderator
	Helper
	MediaManager
	Prime
	Premium
	Media
	Booster
	Partner
	Player
)

var ShortenedRanks = []string{
	"owner", "developer", "manager", "admin", "designer",
	"moderator", "helper", "media_manager", "prime", "premium",
	"media", "booster", "partner", "player",
}

var ChatPrefixes = []string{
	"black", "dark-purple", "blue", "dark-red", "aqua",
	"dark-green", "gold", "lapis", "yellow", "aqua",
	"copper", "purple", "amethyst", "grey",
}

var RankPrefixes = []string{
	"[Owner]", "[Dev]", "[Manager]", "[Admin]", "[Designer]",
	"[Mod]", "[Helper]", "[Media Manager]", "[Prime]", "[Premium]",
	"[Media]", "[Booster]", "[Partner]", "",
}

func (r Rank) Shortened() string {
	return ShortenedRanks[r]
}

func (r Rank) ChatPrefix() string {
	return ChatPrefixes[r]
}

func (r Rank) Prefix() string {
	return text.Colourf("<%v>%v</%v>", r.ChatPrefix(), RankPrefixes[r], r.ChatPrefix())
}

func RankFromPrefix(prefix string) Rank {
	id := Owner
	for _, rank := range RankPrefixes {
		if text.Colourf(rank) == text.Colourf(prefix) {
			return id
		}
		id++
	}

	return -1
}

func RankFromName(name string) Rank {
	id := Owner
	for _, rank := range ShortenedRanks {
		if rank == name {
			return id
		}
		id++
	}

	return -1
}
