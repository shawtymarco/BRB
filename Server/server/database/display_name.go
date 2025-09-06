package database

import (
	"fmt"
	"math"
	"server/server/font"
	"strconv"
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nameConfig struct {
	Rank       bool
	ELO        bool
	TeamColour string
	Health     bool
}

func (nc nameConfig) Name(pd *PlayerData) string {
	var eloStr string
	if nc.ELO {
		eloStr = fmt.Sprintf("<dark-grey>[%v %v]</dark-grey> ", pd.Statistics.ELORank().EloIcon(), font.Transform(strconv.Itoa(pd.Statistics.ELO)))
	}

	var teamColourStr string
	if nc.TeamColour != "" {
		teamColourStr = fmt.Sprintf("<bold><%v>[%v]</%v></bold> ", nc.TeamColour, strings.ToUpper(nc.TeamColour), nc.TeamColour)
	}

	r := pd.Rank()

	n := pd.Username
	if pd.Cosmetics.Nickname != "" {
		n = pd.Cosmetics.Nickname
	}

	if nc.Rank && r != Player {
		return text.Colourf("%v%v%v <%v>%v</%v>", eloStr, teamColourStr, r.Prefix(), r.ChatPrefix(), n, r.ChatPrefix())
	}

	return text.Colourf("%v%v<grey>%v</grey>", eloStr, teamColourStr, n)
}

func (nc nameConfig) NameWithHealth(pd *PlayerData, pl *player.Player) string {
	baseName := nc.Name(pd)

	if !nc.Health {
		return baseName
	}

	health := pl.Health()
	hearts := int(math.Ceil(health))

	if hearts < 0 {
		hearts = 0
	} else if hearts > 20 {
		hearts = 20
	}

	healthStr := fmt.Sprintf("<red>%d❤</red>", hearts)

	return text.Colourf("%v\n%v", baseName, healthStr)
}

var LobbyNameDisplay = nameConfig{Rank: true, ELO: true, Health: false}
var BedWarsNameDisplay = func(teamColour string) nameConfig {
	return nameConfig{TeamColour: teamColour, Health: true}
}
var BuildFFANameDisplay = nameConfig{Health: true}

func TeamColoredName(pd *PlayerData, teamColour string) string {
	n := pd.Username
	if pd.Cosmetics.Nickname != "" {
		n = pd.Cosmetics.Nickname
	}
	return text.Colourf("<%v>%v</%v>", teamColour, n, teamColour)
}
