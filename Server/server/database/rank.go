package database

import "github.com/sandertv/gophertunnel/minecraft/text"

type Rank int

const (
	Owner Rank = iota
	Manager
	Admin
	Moderator
	Helper
	Prime
	Premium
	MediaPartner
	Booster
	Player
)

var ShortenedRanks = []string{
	"owner", "manager", "admin", "moderator", "helper",
	"prime", "premium", "media_partner", "booster", "player",
}

var ChatPrefixes = []string{
	"<lapis>", "<diamond>", "<dark-green>", "<amethyst>", "<copper>",
	"<dark-purple>", "<aqua>", "<gold>", "<quartz>", "<grey>",
}

var RankPrefixes = []string{
	"<lapis>Owner</lapis>", "<diamond>Manager</diamond>", "<dark-green>Admin</dark-green>", "<amethyst>Moderator</amethyst>", "<copper>Helper</copper>",
	"<dark-purple>Prime</dark-purple>", "<aqua>Premium</aqua>", "<gold>Media Partner</gold>", "<quartz>Booster</quartz>", "<grey>Player</grey>",
}

func (r Rank) Shortened() string {
	return ShortenedRanks[r]
}

func (r Rank) ChatPrefix() string {
	return ChatPrefixes[r]
}

func (r Rank) Prefix() string {
	return RankPrefixes[r]
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
