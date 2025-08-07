package database

import (
	"fmt"
	"server/server/font"
	"strconv"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	Rank       bool
	ELO        bool
	TeamColour string
	ShowHealth bool 
}

func (nc nameConfig) Name(pd *PlayerData) string {
	var eloStr string
	if nc.ELO {
		eloStr = fmt.Sprintf("<dark-grey>[%v %v]</dark-grey> ", pd.Statistics.ELORank().EloIcon(), font.Transform(strconv.Itoa(pd.Statistics.ELO)))
	}

	var teamColourStr string
	if nc.TeamColour != "" {
		teamColourStr = fmt.Sprintf("<bold><%v>%v</%v></bold> ", nc.TeamColour, strings.ToUpper(string([]rune(nc.TeamColour)[0])), nc.TeamColour)
	}

	r := pd.Rank()

	n := pd.Username
	if pd.Cosmetics.Nickname != "" {
		n = pd.Cosmetics.Nickname
	}

	var result string
	if nc.Rank && r != Player {
		result = text.Colourf("%v%v%v <%v>%v</%v>", eloStr, teamColourStr, r.Prefix(), r.ChatPrefix(), n, r.ChatPrefix())
	} else {
		result = text.Colourf("%v%v<grey>%v</grey>", eloStr, teamColourStr, n)
	}

	if nc.ShowHealth {
		health := pd.Health 
		result += text.Colourf("\n<red>%v ❤</red>", health)
	}

	return result
}

var LobbyNameDisplay = nameConfig{Rank: true, ELO: true, ShowHealth: false}
var BedWarsNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{TeamColour: teamColour, ShowHealth: true}
}
