package database

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type EloRank int

const (
	Bronze EloRank = iota
	Silver
	Gold
	Platinum
	Diamond
	Emerald
	Sapphire
	Ruby
	Crystal
	Opal
	Amethyst
	Obsidian
	Aventurine
	Quartz
	Topaz
	DarkMatter
)

func GetEloRank(elo int) EloRank {
	if elo < 100 {
		return Bronze
	}

	if elo < 200 {
		return Silver
	}

	if elo < 300 {
		return Gold
	}

	if elo < 400 {
		return Platinum
	}

	if elo < 500 {
		return Diamond
	}

	if elo < 600 {
		return Emerald
	}

	if elo < 700 {
		return Sapphire
	}

	if elo < 800 {
		return Ruby
	}

	if elo < 900 {
		return Crystal
	}

	if elo < 1000 {
		return Opal
	}

	if elo < 1200 {
		return Amethyst
	}

	if elo < 1400 {
		return Obsidian
	}

	if elo < 1600 {
		return Aventurine
	}

	if elo < 1800 {
		return Quartz
	}

	if elo < 2000 {
		return Topaz
	}

	return DarkMatter
}

func (e EloRank) EloPrefix() string {
	switch e {
	case Bronze:
		return text.Colourf("<copper>Bronze</copper>")
	case Silver:
		return text.Colourf("<grey>Silver</grey>")
	case Gold:
		return text.Colourf("<gold>Gold</gold>")
	case Platinum:
		return text.Colourf("<lapis>Platinum</lapis>")
	case Diamond:
		return text.Colourf("<diamond>Diamond</diamond>")
	case Emerald:
		return text.Colourf("<emerald>Emerald</emerald>")
	case Sapphire:
		return text.Colourf("<dark-blue>Sapphire</dark-blue>")
	case Ruby:
		return text.Colourf("<redstone>Ruby</redstone>")
	case Crystal:
		return text.Colourf("<iron>Crystal</iron>")
	case Opal:
		return text.Colourf("<white>Opal</white>")
	case Amethyst:
		return text.Colourf("<amethyst>Amethyst</amethyst>")
	case Obsidian:
		return text.Colourf("<dark-purple>Obsidian</dark-purple>")
	case Aventurine:
		return text.Colourf("<green>Aventurine</green>")
	case Quartz:
		return text.Colourf("<quartz>Quartz</quartz>")
	case Topaz:
		return text.Colourf("<yellow>Topaz</yellow>")
	case DarkMatter:
		return text.Colourf("<black>Dark Matter</black>")
	default:
		return ""
	}
}

func (e EloRank) EloIcon() string {
	switch e {
	case Bronze:
		return "\uE000"
	case Silver:
		return "\uE001"
	case Gold:
		return "\uE002"
	case Platinum:
		return "\uE003"
	case Diamond:
		return "\uE004"
	case Emerald:
		return "\uE005"
	case Sapphire:
		return "\uE006"
	case Ruby:
		return "\uE007"
	case Crystal:
		return "\uE008"
	case Opal:
		return "\uE009"
	case Amethyst:
		return "\uE00A"
	case Obsidian:
		return "\uE00B"
	case Aventurine:
		return "\uE00C"
	case Quartz:
		return "\uE00D"
	case Topaz:
		return "\uE00E"
	case DarkMatter:
		return "\uE00F"
	}

	return ""
}
