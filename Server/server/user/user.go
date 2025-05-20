package user

import (
	"errors"
	"fmt"
	"math/rand"
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/utils"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	userMu sync.Mutex
	users  = map[uuid.UUID]*User{}
)

type User struct {
	pl         *player.Player
	h          *world.EntityHandle
	Scoreboard *scoreboard.Scoreboard
	Data       *database.PlayerData
	FirstTime  bool
}

func New(pl *player.Player, isBot bool) (*User, error) {
	if pl == nil {
		return nil, fmt.Errorf("new player should not be nil")
	}

	userMu.Lock()
	defer userMu.Unlock()
	ft := false
	if DataFromPlayer(pl) == nil {
		ft = true
		pd := &database.PlayerData{
			Uuid:      pl.UUID(),
			Username:  pl.Name(),
			FirstJoin: time.Now(),
			LastJoin:  time.Now(),
			GroupSettings: database.GroupSettings{
				RankId: database.Player.Shortened(),
			},
			GlobalStats: database.GlobalStats{
				Level: 1,
			},
			Games: database.Games{},
		}
		if err := server.Database.CreatePlayer(pd); err != nil {
			return nil, err
		}
	}

	d, err := server.Database.FindPlayer(pl.UUID())
	if err != nil {
		return nil, err
	}

	if pl.Name() != d.Username {
		d.Username = pl.Name()
	}

	d.Online = true
	if !isBot {
		d.ProtocolId = utils.Session(pl).ClientData().GameVersion
	}

	u := &User{
		pl:        pl,
		h:         pl.H(),
		Data:      d,
		FirstTime: ft,
	}
	users[pl.UUID()] = u
	return u, nil
}

func Save(pl *player.Player) {
	user := LookupUUID(pl.UUID())
	user.Data.LastJoin = time.Now()
	user.Data.Online = false
	utils.Panic(server.Database.SavePlayer(user.Data))
}

func Delete(pl *player.Player) {
	userMu.Lock()
	defer userMu.Unlock()
	delete(users, pl.UUID())
}

func LookupPlayer(pl *player.Player) *User {
	return LookupUUID(pl.UUID())
}

func LookupUUID(uuid uuid.UUID) *User {
	userMu.Lock()
	defer userMu.Unlock()
	return users[uuid]
}

func DataFromPlayer(pl *player.Player) *database.PlayerData {
	d, err := server.Database.FindPlayer(pl.UUID())
	if err != nil {
		var playerDataNotFoundError utils.PlayerDataNotFoundError
		if !errors.As(err, &playerDataNotFoundError) {
			utils.Panic(err)
		}
		return nil
	}
	return d
}

func (u *User) Player() *player.Player {
	return u.pl
}

func (u *User) H() *world.EntityHandle {
	return u.h
}

func (u *User) AddItem(its ...item.Stack) bool {
	if len(u.pl.Inventory().Items())+len(its) > 36 {
		u.pl.Message(text.Colourf(language.Translate(u.pl).Global.Error.InventoryFull))
		return false
	}
	for _, it := range its {
		_, err := u.pl.Inventory().AddItem(it)
		if err != nil {
			u.pl.Message(text.Colourf(language.Translate(u.pl).Global.Error.InventoryFull))
			return false
		}
	}
	return true
}

func (u *User) DropItem(it item.Stack, tx *world.Tx) {
	ent := entity.NewItem(world.EntitySpawnOpts{
		Position: u.pl.Position(),
		Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1},
	}, it)
	tx.AddEntity(ent)
}

func (u *User) PlaySound(sound CustomSound, title, author string, volume, pitch float64) {
	pos := u.pl.Position()
	pk := &packet.PlaySound{
		SoundName: string(sound),
		Position:  mgl32.Vec3{float32(pos.X()), float32(pos.Y()), float32(pos.Z())},
		Volume:    float32(volume),
		Pitch:     float32(pitch),
	}
	utils.WritePacket(utils.Session(u.pl), pk)
	if title != "" && author != "" {
		u.pl.Message(text.Colourf(language.Translate(u.pl).Global.Misc.NowPlaying, server.Config.Prefix, title, author))
	}
}

func SetSpectator(pl *player.Player, set bool) {
	if !set {
		pl.SetGameMode(world.GameModeSurvival)
		return
	}

	pl.SetGameMode(SpectatorGamemode{})
	pl.SetFlightSpeed(0.15)
}
