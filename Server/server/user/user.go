package user

import (
	"errors"
	"fmt"
	"math/rand"
	"server/server"
	"server/server/cooldown"
	"server/server/database"
	"server/server/font"
	"server/server/game"
	"server/server/language"
	"server/server/utils"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block"

	skin2 "github.com/df-mc/dragonfly/server/player/skin"

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
	pl *player.Player
	h  *world.EntityHandle

	cooldownMap cooldown.MappedCoolDown[PlayerCoolDowns]

	Scoreboard *scoreboard.Scoreboard
	Data       *database.PlayerData
	FirstTime  bool

	GameInfo GameRuntimeData
	Game     *game.Game

	LastHit   *world.EntityHandle
	LastHitAt time.Time
}

func New(pl *player.Player, isBot bool) (*User, error) {
	if pl == nil {
		return nil, fmt.Errorf("new player should not be nil")
	}

	userMu.Lock()
	defer userMu.Unlock()

	ft := false
	if !isBot && DataFromPlayer(pl) == nil {
		ft = true
		pd := &database.PlayerData{
			Uuid:      pl.UUID(),
			Username:  pl.Name(),
			FirstJoin: time.Now(),
			LastJoin:  time.Now(),
			DeviceOS:  utils.Session(pl).ClientData().DeviceOS,
			Statistics: database.Statistics{
				RankId: database.Player.Shortened(),
				Level:  1,
			},
			Cosmetics: database.Cosmetics{
				SelectedWoodType: block.OakWood(),
			},
			Settings: database.Settings{
				HotBarConfig:   [9]database.HotBarCategory(make([]database.HotBarCategory, 9)),
				QuickBuyConfig: make(map[int]*int),
			},
		}
		if err := server.Database.CreatePlayer(pd); err != nil {
			return nil, err
		}
	}

	var d *database.PlayerData
	if isBot {
		d = utils.Panics(server.Database.FindPlayerFromName(pl.Name(), &database.PlayerNameSearchOpts{CaseInsensitive: false, PartialMatch: false}))
	} else {
		d = utils.Panics(server.Database.FindPlayer(pl.UUID()))
	}

	if pl.Name() != d.Username {
		d.Username = pl.Name()
	}

	d.Online = true
	if !isBot {
		d.ProtocolId = utils.Session(pl).ClientData().GameVersion
	}

	u := &User{
		pl: pl,
		h:  pl.H(),

		cooldownMap: cooldown.NewMappedCoolDown[PlayerCoolDowns](),

		Data:      d,
		FirstTime: ft,
	}
	users[pl.UUID()] = u
	return u, nil
}

func Save(pl *player.Player) {
	user := LookupPlayer(pl)
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

func LookupUserID(userId string) *User {
	userMu.Lock()
	defer userMu.Unlock()

	for _, user := range users {
		if user.Data.UserId == userId {
			return user
		}
	}
	return nil
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

// IsCooldownActive sets a cooldown for the given type if it doesn't already exist or renews it if specified.
// It returns true if the cooldown is initially active (to be used for conditional command execution).
func (u *User) IsCooldownActive(cooldownType PlayerCoolDowns, duration time.Duration, renew, sendMessage bool) bool {
	coolDown := u.cooldownMap[cooldownType]

	if coolDown == nil {
		coolDown = cooldown.NewCoolDown()
		u.cooldownMap[cooldownType] = coolDown
	}
	exists := coolDown.Active()
	if renew || !coolDown.Active() {
		coolDown.Set(duration)
	}

	if sendMessage && exists {
		u.pl.Message(text.Colourf(language.Translate(u.pl).Commands.Error.CoolDown, coolDown.Remaining()))
	}

	return exists
}

func (u *User) AddItemWithHBConfig(preferredSlot int, it item.Stack) (n int, err error) {
	var category database.HotBarCategory

	switch it.Item().(type) {
	case item.Sword, item.Stick:
		category = database.Melee
	case world.Block:
		category = database.Blocks
	case item.Bow:
		category = database.Bows
	case item.Potion:
		category = database.Potions
	case item.GoldenApple, item.Snowball, block.TNT, item.EnderPearl, item.Bucket, item.Egg, block.Sponge, item.Compass:
		category = database.Utility
	case item.Shears:
		category = database.Shears
	case item.Pickaxe:
		category = database.Pickaxe
	case item.Axe:
		category = database.Axe
	}

	slot := -1

	if category != database.None {
		for s, c := range u.Data.Settings.HotBarConfig {
			if c == category {
				slot = s
				break
			}
		}
	}

	if slot != -1 {
		sa, sb := it.AddStack(utils.Panics(u.pl.Inventory().Item(slot)))
		if sb.Count() != 0 {
			slot = -1
		} else {
			return it.Count(), u.pl.Inventory().SetItem(slot, sa)
		}
	}

	if slot == -1 && preferredSlot != -1 {
		sa, sb := it.AddStack(utils.Panics(u.pl.Inventory().Item(preferredSlot)))
		if sb.Count() == 0 {
			return it.Count(), u.pl.Inventory().SetItem(preferredSlot, sa)
		}
	}

	return u.pl.Inventory().AddItem(it)
}

func (u *User) SendScoreboard(numOfSpaces int) {
	u.Scoreboard.Set(0, text.Colourf("%v%v", strings.Repeat(" ", numOfSpaces), font.Transform("SEASON 1")))
	u.pl.SendScoreboard(u.Scoreboard)
}

func (u *User) DropItem(it item.Stack, tx *world.Tx) {
	ent := entity.NewItem(world.EntitySpawnOpts{
		Position: u.pl.Position(),
		Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1},
	}, it)
	tx.AddEntity(ent)
}

func (u *User) PlaySound(sound string, title, author string, volume, pitch float64) {
	pos := u.pl.Position()
	pk := &packet.PlaySound{
		SoundName: sound,
		Position:  mgl32.Vec3{float32(pos.X()), float32(pos.Y()), float32(pos.Z())},
		Volume:    float32(volume),
		Pitch:     float32(pitch),
	}
	utils.WritePacket(utils.Session(u.pl), pk)
	if title != "" && author != "" {
		u.pl.Message(text.Colourf(language.Translate(u.pl).Misc.NowPlaying, server.Config.Prefix, title, author))
	}
}

func (u *User) RefreshCape() {
	cape, ok := database.GetCapeByType(u.Data.Cosmetics.SelectedCape)
	if !ok {
		return
	}

	skin := u.pl.Skin()
	skin.Cape = skin2.NewCape(64, 32)
	skin.Cape.Pix = database.CapeAsBytes(cape)

	go func() {
		u.h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			e.(*player.Player).SetSkin(skin)
		})
	}()
}

func (u *User) SetSpectator(set bool) {
	if !set {
		u.pl.SetGameMode(world.GameModeSurvival)
		u.Game.RemoveSpectator(u.pl)
		return
	}

	u.pl.SetGameMode(SpectatorGamemode{})
	u.pl.SetFlightSpeed(0.15)
	u.Game.AddSpectator(u.pl)
}

type GameRuntimeData struct {
	BedWars struct {
		Kills      int
		FinalKills int
		BedsBroken int
	}
	BuildFFA struct {
		Kills int
	}
}

func (d GameRuntimeData) TotalBWKills() int {
	return d.BedWars.Kills + d.BedWars.FinalKills
}
