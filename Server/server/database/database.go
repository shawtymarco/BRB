package database

import "github.com/google/uuid"

type Database interface {
	String() string
	CreatePlayer(data *PlayerData) error
	DeletePlayerByName(playerName string, opts *PlayerNameSearchOpts) error
	SavePlayer(data *PlayerData) error
	FindPlayer(uuid uuid.UUID) (*PlayerData, error)
	FindPlayerByDiscordID(id string) (*PlayerData, error)
	FindPlayerByName(playerName string, opts *PlayerNameSearchOpts) (*PlayerData, error)
	FindAllPlayers() ([]*PlayerData, error)
	SaveAll() map[string]error
}

type PlayerNameSearchOpts struct {
	CaseInsensitive bool
	PartialMatch    bool
}
