package database

import (
	"fmt"
	"regexp"
	"server/server/utils"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func NewLocalDatabase() *LocalDatabase {
	return &LocalDatabase{data: make(map[uuid.UUID]*PlayerData)}
}

type LocalDatabase struct {
	data map[uuid.UUID]*PlayerData
	mu   sync.RWMutex
}

func (d *LocalDatabase) String() string {
	return "Local"
}

func (d *LocalDatabase) CreatePlayer(data *PlayerData) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[data.Uuid] = data
	return nil
}

func (d *LocalDatabase) DeletePlayerByName(playerName string, opts *PlayerNameSearchOpts) error {
	player, err := d.FindPlayerByName(playerName, opts)
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.data, player.Uuid)
	return nil
}

func (d *LocalDatabase) SavePlayer(data *PlayerData) error {
	return d.CreatePlayer(data)
}

func (d *LocalDatabase) FindPlayer(uuid uuid.UUID) (*PlayerData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.data[uuid] != nil {
		return d.data[uuid], nil
	}
	return nil, utils.PlayerDataNotFoundError{Identifier: uuid.String()}
}

func (d *LocalDatabase) FindPlayerByDiscordID(id string) (*PlayerData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, playerData := range d.data {
		if playerData.UserId == id {
			return playerData, nil
		}
	}
	return nil, utils.PlayerDataNotFoundError{Identifier: id}
}

func (d *LocalDatabase) FindPlayerByName(playerName string, opts *PlayerNameSearchOpts) (*PlayerData, error) {
	if opts == nil {
		opts = &PlayerNameSearchOpts{}
	}
	patternStr := playerName
	if opts.CaseInsensitive {
		patternStr = strings.ToLower(patternStr)
	}
	if !opts.PartialMatch {
		patternStr = fmt.Sprintf("^%v$", patternStr)
	}
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, playerData := range d.data {
		for _, name := range playerData.AlternativeMCAccounts {
			if opts.CaseInsensitive {
				name = strings.ToLower(name)
			}
			if pattern.MatchString(name) {
				return playerData, nil
			}
		}
	}
	return nil, utils.PlayerDataNotFoundError{Identifier: playerName}
}

func (d *LocalDatabase) SaveAll() map[string]error {
	return make(map[string]error) // noop
}
