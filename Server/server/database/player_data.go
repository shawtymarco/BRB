package database

import (
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PlayerData struct {
	Uuid   uuid.UUID
	UserId string

	AlternativeMCAccounts []string
	AlternativeDCAccounts []string

	Username      string
	Online        bool
	FirstJoin     time.Time
	LastJoin      time.Time
	ProtocolId    string
	GroupSettings GroupSettings
	GlobalStats   GlobalStats
	Games         Games
}

type GroupSettings struct {
	RankId      string
	Permissions []string
}

func (gs GroupSettings) Rank() Rank {
	return RankFromName(gs.RankId)
}

type GlobalStats struct {
	Level int
	Xp    int
	RR    int
	Coins int
}

func (gs GlobalStats) LevelDisplay() string {
	return text.Colourf(Lvls[gs.Level/10](strconv.Itoa(gs.Level)))
}

func (gs GlobalStats) MaxXP() int {
	if gs.Level == 1 {
		return 100 // Base XP for level 1
	} else if gs.Level >= 100 {
		return -1 // Cap at level 100
	} else {
		growthFactor := 1.05
		return int(100 * (math.Pow(growthFactor, float64(gs.Level-1))))
	}
}

type Games struct {
	TowerWars TowerWars
	Pvp       Pvp
}

type TowerWars struct {
	Wins         int
	Losses       int
	TowersPlaced int
	UnitsKilled  int
}

type Pvp struct {
	TotalKills       int
	TotalKillStreaks int
	TotalDeaths      int
	Duels            Duels
	FFA              FFA
}

func (p Pvp) KillDeathRatio() float64 {
	x := float64(p.TotalKills) / float64(p.TotalDeaths)
	return math.Ceil(x*10) / 10
}

type Duels struct {
	Magic    KDData
	BedFight KDData
	Pearl    KDData
	Sumo     KDData
	Bot      KDData
}

type FFA struct {
	Sumo       KDData
	Build      KDData
	Resistance Resistance
}

type KDData struct {
	Kills      int
	Deaths     int
	KillStreak int
}

type Resistance struct {
	Hits int
}
