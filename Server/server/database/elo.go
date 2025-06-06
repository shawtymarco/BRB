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
		return text.Colourf("<bold><copper>Bronze</copper></bold>")
	case Silver:
		return text.Colourf("<bold><grey>Silver</grey></bold>")
	case Gold:
		return text.Colourf("<bold><gold>Gold</gold></bold>")
	case Platinum:
		return text.Colourf("<bold><lapis>Platinum</lapis></bold>")
	case Diamond:
		return text.Colourf("<bold><diamond>Diamond</diamond></bold>")
	case Emerald:
		return text.Colourf("<bold><emerald>Emerald</emerald></bold>")
	case Sapphire:
		return text.Colourf("<bold><dark-blue>Sapphire</dark-blue></bold>")
	case Ruby:
		return text.Colourf("<bold><redstone>Ruby</redstone></bold>")
	case Crystal:
		return text.Colourf("<bold><iron>Crystal</iron></bold>")
	case Opal:
		return text.Colourf("<bold><white>Opal</white></bold>")
	case Amethyst:
		return text.Colourf("<bold><amethyst>Amethyst</amethyst></bold>")
	case Obsidian:
		return text.Colourf("<bold><dark-purple>Obsidian</dark-purple></bold>")
	case Aventurine:
		return text.Colourf("<bold><green>Aventurine</green></bold>")
	case Quartz:
		return text.Colourf("<bold><quartz>Quartz</quartz></bold>")
	case Topaz:
		return text.Colourf("<bold><yellow>Topaz</yellow></bold>")
	case DarkMatter:
		return text.Colourf("<bold><black>Dark Matter</black></bold>")
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
