package database

import (
	"fmt"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	Level      bool
	TeamColour string
}

func (nc nameConfig) Name(pd *PlayerData) string {
	var levelStr string
	if nc.Level {
		l := pd.GlobalStats.Level
		levelStr = fmt.Sprintf("<dark-grey>[%v]</dark-grey> ", l)
	}

	var teamColourStr string
	if nc.TeamColour != "" {
		teamColourStr = fmt.Sprintf("<%v>[%v]</%v> ", nc.TeamColour, strings.ToUpper(nc.TeamColour), nc.TeamColour)
	}

	n := pd.Username

	if pd.GroupSettings.Rank() == Player {
		return text.Colourf("%v%v%v%v", levelStr, teamColourStr, pd.GroupSettings.Rank().ChatPrefix(), n)
	}

	r := pd.GroupSettings.Rank()
	return text.Colourf("%v%v%v %v", teamColourStr, levelStr, r.ChatPrefix(), n)
}

var BasicNameDisplay = nameConfig{}
var LobbyNameDisplay = nameConfig{Level: true}
var BedWarsNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{Level: true, TeamColour: teamColour}
}
