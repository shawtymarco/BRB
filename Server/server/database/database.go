package database

import "github.com/google/uuid"

type Database interface {
	// String returns the name of this database type.
	String() string
	// CreatePlayer saves the new *unique* player data.
	CreatePlayer(data *PlayerData) error
	// DeletePlayerByName deletes the player data using a player name
	DeletePlayerByName(playerName string, opts *PlayerNameSearchOpts) error
	// SavePlayer overwrites the existing player data with this new data.
	SavePlayer(data *PlayerData) error
	// FindPlayer finds a player from uuid.
	FindPlayer(uuid uuid.UUID) (*PlayerData, error)
	// FindPlayerByDiscordID finds a player using their Discord's ID (must be registered with Discord)
	FindPlayerByDiscordID(id string) (*PlayerData, error)
	// FindPlayerByName finds a player from a player name.
	FindPlayerByName(playerName string, opts *PlayerNameSearchOpts) (*PlayerData, error)
	// SaveAll saves all loaded players to the database.
	SaveAll() map[string]error
}

type PlayerNameSearchOpts struct {
	CaseInsensitive bool
	PartialMatch    bool
}
