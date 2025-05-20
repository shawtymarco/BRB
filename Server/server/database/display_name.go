package database

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	TeamColour string
	Level      bool
}

func (nc nameConfig) Name(pd *PlayerData) string {
	var levelStr string
	if nc.Level {
		l := pd.GlobalStats.LevelDisplay()
		levelStr = fmt.Sprintf("<dark-grey>[<dark-grey>%v</dark-grey>]</dark-grey> ", l)
	}

	var teamColourStr string
	if nc.TeamColour != "" {
		teamColourStr = fmt.Sprintf("<dark-grey>[<dark-grey>%v</dark-grey>]</dark-grey> ", nc.TeamColour)
	}

	n := pd.Username

	if pd.GroupSettings.Rank().Icon() == "" {
		return text.Colourf("%v%v%v%v", teamColourStr, levelStr, pd.GroupSettings.Rank().ChatPrefix(), n)
	}

	r := pd.GroupSettings.Rank()
	return text.Colourf("%v%v%v%v %v", teamColourStr, levelStr, r.Icon(), r.ChatPrefix(), n)
}

var BasicNameDisplay = nameConfig{}
var LobbyNameDisplay = nameConfig{Level: true}
var TWPlayingNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{Level: true, TeamColour: teamColour}
}

var Lvls = []func(lvl string) string{
	func(lvl string) string {
		return "<quartz>" + lvl + "</quartz>"
	},
	func(lvl string) string {
		return "<copper>" + lvl + "</copper>"
	},
	func(lvl string) string {
		return "<iron>" + lvl + "</iron>"
	},
	func(lvl string) string {
		return "<lapis>" + lvl + "</lapis>"
	},
	func(lvl string) string {
		return "<redstone>" + lvl + "</redstone>"
	},
	func(lvl string) string {
		return "<bold><italic><gold>" + lvl + "</gold></italic></bold>"
	},
	func(lvl string) string {
		return "<bold><italic><amethyst>" + lvl + "</amethyst></italic></bold>"
	},
	func(lvl string) string {
		return "<bold><italic><emerald>" + lvl + "</emerald></italic></bold>"
	},
	func(lvl string) string {
		return "<bold><italic><diamond>" + lvl + "</diamond></italic></bold>"
	},
	func(lvl string) string {
		return "<obfuscated><black>kk</black></obfuscated> <bold><italic><netherite>" + lvl + "</netherite></italic></bold> <obfuscated><black>kk</black></obfuscated>"
	},
}
