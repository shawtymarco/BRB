package database

import (
	"fmt"
	"github.com/samber/lo"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	Rank       bool
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
		teamColourStr = fmt.Sprintf("<bold><%v>%v</%v></bold> ", nc.TeamColour, strings.ToUpper(string([]rune(nc.TeamColour)[0])), nc.TeamColour)
	}

	n := pd.Username

	if pd.Rank() == Player {
		return text.Colourf("%v%v%v%v", eloStr, teamColourStr, pd.Rank().ChatPrefix(), n)
	}

	r := pd.Rank()
	return text.Colourf("%v%v%v %v", teamColourStr, eloStr, lo.If(nc.Rank, r.ChatPrefix()).Else(""), n)
}

var LobbyNameDisplay = nameConfig{Rank: true, ELO: true}
var BedWarsNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{TeamColour: teamColour}
}
