package database

import (
	"fmt"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	ELO        bool
	TeamColour string
}

func (nc nameConfig) Name(pd *PlayerData) string {
	var eloStr string
	if nc.ELO {
		eloStr = fmt.Sprintf("<dark-grey>[%v %v]</dark-grey> ", pd.Statistics.ELORank().EloIcon(), pd.Statistics.ELO)
	}

	var teamColourStr string
	if nc.TeamColour != "" {
		teamColourStr = fmt.Sprintf("<%v>[%v]</%v> ", nc.TeamColour, strings.ToUpper(nc.TeamColour), nc.TeamColour)
	}

	n := pd.Username

	if pd.Rank() == Player {
		return text.Colourf("%v%v%v%v", eloStr, teamColourStr, pd.Rank().ChatPrefix(), n)
	}

	r := pd.Rank()
	return text.Colourf("%v%v%v %v", teamColourStr, eloStr, r.ChatPrefix(), n)
}

var LobbyNameDisplay = nameConfig{ELO: true}
var BedWarsNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{ELO: true, TeamColour: teamColour}
}
