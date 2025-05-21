package database

import "github.com/google/uuid"

type Database interface {
	// String returns the name of this database type.
	String() string
	// CreatePlayer saves the new *unique* player data.
	CreatePlayer(data *PlayerData) error
	// SavePlayer overwrites the existing player data with this new data.
	SavePlayer(data *PlayerData) error
	// FindPlayer finds a player from a uuid.
	FindPlayer(uuid uuid.UUID) (*PlayerData, error)
	// FindPlayerFromName finds a player from a player name.
	FindPlayerFromName(playerName string, opts *PlayerNameSearchOpts) (*PlayerData, error)
	// SaveAll saves all loaded players to the database.
	SaveAll() map[string]error
}

type PlayerNameSearchOpts struct {
	CaseInsensitive bool
	PartialMatch    bool
}
