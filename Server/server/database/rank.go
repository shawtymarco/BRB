package database

import "github.com/sandertv/gophertunnel/minecraft/text"

type Rank int

const (
	Owner Rank = iota
	Council
	Sheriff
	Mason
	Developer
	Artist
	Partner
	Platinum
	Gold
	Silver
	VIP
	Player
)

var RankIcons = []string{
	"\uE000", "\uE001", "", "\uE002", "\uE003",
	"\uE004", "\uE005", "\uE006", "\uE007", "\uE008",
	"\uE009", "",
}

var ShortenedRanks = []string{
	"owner", "council", "sheriff", "mason", "developer",
	"artist", "partner", "platinum", "gold", "silver",
	"vip", "player",
}

var ChatPrefixes = []string{
	"<lapis>", "<diamond>", "<dark-green>", "<amethyst>", "<amethyst>",
	"<copper", "<dark-purple>", "<aqua>", "<gold>", "<quartz>",
	"<emerald>", "<grey>",
}

var RankPrefixes = []string{
	"<lapis>Owner</lapis>", "<diamond>Council</diamond>", "<dark-green>Sheriff</dark-green>", "<amethyst>Mason</amethyst>", "<amethyst>Developer</amethyst>",
	"<copper>Artist</copper>", "<dark-purple>Partner</dark-purple>", "<aqua>Platinum</aqua>", "<gold>Gold</gold>", "<quartz>Silver</quartz>",
	"<emerald>VIP</emerald>", "<grey>Player</grey>",
}

func (r Rank) Icon() string {
	return RankIcons[r]
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
