package worldmanager

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/world"
	"golang.org/x/exp/maps"
)

var Worlds []string

type ManagerSettings struct {
	Folder   string
	Fallback *world.World
	Logger   *slog.Logger
}

func (s ManagerSettings) NewManager() (*Manager, error) {
	possibleWorlds := make(map[string]WorldEntry)
	worldFolders, err := os.ReadDir(s.Folder)
	if err != nil {
		return nil, err
	}
	for _, dir := range worldFolders {
		cleanName := dir.Name()
		fullPath := fmt.Sprintf("%v/%v", s.Folder, cleanName)
		zipped := false
		if !dir.IsDir() && path.Ext(dir.Name()) == ".mcworld" {
			cleanName = strings.Split(dir.Name(), ".")[0]
			fullPath, err = unzipFile(fullPath)
			if err != nil {
				return nil, err
			}
			zipped = true
		}
		Worlds = append(Worlds, cleanName)
		possibleWorlds[cleanName] = WorldEntry{
			WorldPath: fullPath,
			WasZipped: zipped,
		}
	}

	return &Manager{log: s.Logger, fallback: s.Fallback, possibleWorlds: possibleWorlds, loaded: make(map[string]*world.World)}, nil
}

var errNotFound = errors.New("world not found")

type WorldEntry struct {
	WorldPath string
	WasZipped bool
}

type Manager struct {
	log      *slog.Logger
	fallback *world.World

	possibleWorlds map[string]WorldEntry
	loaded         map[string]*world.World
}

func (m *Manager) Worlds() []string {
	return maps.Keys(m.possibleWorlds)
}

func (m *Manager) LoadedWorlds() []string {
	return maps.Keys(m.loaded)
}

func (m *Manager) World(worldName string, new bool) (*world.World, error) {
	if w, ok := m.loaded[worldName]; ok && new {
		return w, nil
	}
	worldEntry, ok := m.possibleWorlds[worldName]
	if !ok {
		return m.fallback, errNotFound
	}
	w, err := World(worldEntry.WorldPath, false, m.log)
	if err != nil {
		return m.fallback, err
	}
	m.loaded[worldName] = w
	return w, nil
}

func (m *Manager) UnloadWorld(worldName string) error {
	w, ok := m.loaded[worldName]
	if !ok {
		return errNotFound
	}
	if m.fallback != nil {
		w.Exec(func(tx *world.Tx) {
			for entity := range tx.Entities() {
				m.fallback.Exec(func(tx1 *world.Tx) {
					tx1.AddEntity(entity.H())
				})
			}
		})
	}
	delete(m.loaded, worldName)
	if err := w.Close(); err != nil {
		return err
	}
	if worldEntry := m.possibleWorlds[worldName]; worldEntry.WasZipped {
		zipPath, err := zipDir(worldEntry.WorldPath)
		if err != nil {
			return err
		}
		if err = os.Rename(zipPath, fmt.Sprintf("./%v.%v.zip", worldName, time.Now().Format("2006-01-02.15-04-05"))); err != nil {
			return err
		}
	}

	return nil
}
