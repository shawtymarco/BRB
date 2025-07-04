package tebex

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// ExecuteCommands queries the Tebex API for a list of online commands from a *player.Player. If any are returned, the
// commands will automatically be executed.
func (c *Client) ExecuteCommands(p *player.Player) {
	var result struct {
		Commands []struct {
			ID         int    `json:"id"`
			Command    string `json:"command"`
			Conditions struct {
				Delay int `json:"delay,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"conditions"`
		} `json:"commands"`
	}
	err := c.get("queue/online-commands/"+p.XUID(), &result)
	if err != nil {
		c.log.Error("failed to query pending online commands: %v", err)
		return
	}

	var ids []int
	inv := p.Inventory()
	for _, entry := range result.Commands {
		if entry.Conditions.Slots > inv.Size()-len(inv.Items()) {
			// We don't have enough slots to execute this command.
			continue
		}

		time.AfterFunc(time.Second*time.Duration(entry.Conditions.Delay), func() {
			err = c.processCommand(entry.Command, p.Name(), p.XUID())
			if err != nil {
				c.log.Error("failed to execute online command: %s", err)
			}
			err = c.delete("queue", map[string]any{"ids": []int{entry.ID}})
			if err != nil {
				c.log.Error("failed to delete processed online command: %s", err)
			}
		})
		if entry.Conditions.Delay == 0 {
			ids = append(ids, entry.ID)
		}
	}
	if len(ids) > 0 {
		if err = c.delete("queue", map[string]any{"ids": ids}); err != nil {
			c.log.Error("failed to delete processed online commands: %s", err)
		}
	}
}

// ExecuteOfflineCommands queries the Tebex API for a list of offline commands. If any are returned, the commands will
// automatically be executed.
func (c *Client) ExecuteOfflineCommands() {
	var result struct {
		Commands []struct {
			ID         int    `json:"id"`
			Command    string `json:"command"`
			Conditions struct {
				Delay int `json:"delay,omitempty"`
				Slots int `json:"slots,omitempty"`
			} `json:"conditions"`
			Player struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				UUID string `json:"uuid"`
			} `json:"player"`
		} `json:"commands"`
	}
	err := c.get("queue/offline-commands", &result)
	if err != nil {
		c.log.Error("failed to query pending offline commands: %v", err)
		return
	}

	var ids []int
	for _, entry := range result.Commands {
		if entry.Conditions.Delay > 0 {
			time.AfterFunc(time.Second*time.Duration(entry.Conditions.Delay), func() {
				err = c.processCommand(entry.Command, entry.Player.Name, entry.Player.UUID)
				if err != nil {
					c.log.Error("failed to execute offline command: %s", err)
				}
				err = c.delete("queue", map[string]any{"ids": []int{entry.ID}})
				if err != nil {
					c.log.Error("failed to delete processed offline command: %s", err)
				}
				fmt.Println("deleted")
			})
			continue
		}
		err = c.processCommand(entry.Command, entry.Player.Name, entry.Player.UUID)
		if err != nil {
			c.log.Error("failed to execute offline command: %s", err)
		}
		ids = append(ids, entry.ID)
	}
	if len(ids) > 0 {
		if err = c.delete("queue", map[string]any{"ids": ids}); err != nil {
			c.log.Error("failed to delete processed offline commands: %s", err)
		}
	}
}

// Source is a dummy source to be used for command execution.
type Source struct {
	log *slog.Logger
}

// Name ...
func (Source) Name() string {
	return "Tebex"
}

// Position ...
func (Source) Position() mgl64.Vec3 {
	return mgl64.Vec3{}
}

// World ...
func (Source) World() *world.World {
	return nil
}

// SendCommandOutput ...
func (d Source) SendCommandOutput(output *cmd.Output) {
	for _, e := range output.Errors() {
		d.log.Error("error whilst executing commands: %v", e)
	}
}

// processCommand processes the execution of a command line.
func (c *Client) processCommand(commandLine, player, id string) error {
	name := strings.TrimPrefix(strings.Split(commandLine, " ")[0], "/")
	command, ok := cmd.ByAlias(name)
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
	}

	command.Execute(strings.TrimPrefix(strings.TrimPrefix(strings.ReplaceAll(
		strings.ReplaceAll(
			commandLine,
			"{id}",
			id,
		),
		"{username}",
		player,
	), name), " "), Source{log: c.log}, nil)
	return nil
}
