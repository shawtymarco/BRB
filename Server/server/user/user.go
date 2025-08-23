package user

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"server/server"
	"server/server/cooldown"
	"server/server/database"
	"server/server/font"
	"server/server/game"
	"server/server/language"
	"server/server/utils"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/session"

	"github.com/sandertv/gophertunnel/minecraft/protocol"

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

	LastMessenger *world.EntityHandle

	PendingDuelRequestTo *world.EntityHandle

	OldArmour OldArmour

	BuyMutex sync.Mutex
}

func newUser(pl *player.Player, isBot bool) (*User, error) {
	userMu.Lock()
	defer userMu.Unlock()

	ft := false
	var d *database.PlayerData
	if isBot {
		d, _ = server.Database.FindPlayerByName(pl.Name(), &database.PlayerNameSearchOpts{CaseInsensitive: false, PartialMatch: false})
	} else {
		d, _ = server.Database.FindPlayer(pl.UUID())
	}
	if d == nil {
		ft = true
		pd := &database.PlayerData{
			UUID:      pl.UUID(),
			Username:  pl.Name(),
			FirstJoin: time.Now(),
			LastJoin:  time.Now(),
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

		if isBot {
			pd.DeviceOS = protocol.DeviceAndroid
		} else {
			pd.DeviceOS = utils.Session(pl).ClientData().DeviceOS
		}

		if err := server.Database.CreatePlayer(pd); err != nil {
			return nil, err
		}
		d = pd
	}

	if pl.Name() != d.Username {
		d.Username = pl.Name()
	}

	d.Online = true
	if !isBot {
		d.ProtocolId = utils.Session(pl).ClientData().GameVersion
		d.DeviceID = utils.Session(pl).ClientData().DeviceID

		conn := utils.FetchPrivateField[session.Conn](utils.Session(pl), "conn")
		hasher := sha256.New()
		hasher.Write([]byte(conn.RemoteAddr().String() + "big fat boobs"))
		d.HashedIP = hex.EncodeToString(hasher.Sum(nil))
		d.IPStoredSince = time.Now()
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

func GetUser(pl *player.Player) *User {
	if users[pl.UUID()] == nil {
		return utils.Panics(newUser(pl, slices.Contains([]string{"Mark", "Sam", "Steven"}, pl.Name())))
	}

	userMu.Lock()
	defer userMu.Unlock()

	users[pl.UUID()].pl = pl
	users[pl.UUID()].h = pl.H()

	return users[pl.UUID()]
}

func GetUserByUserID(userId string) *User {
	userMu.Lock()
	defer userMu.Unlock()

	for _, user := range users {
		if user.Data.UserId == userId {
			return user
		}
	}
	return nil
}

func Save(pl *player.Player) {
	user := GetUser(pl)
	user.Data.LastJoin = time.Now()
	user.Data.Online = false
	_ = server.Database.SavePlayer(user.Data)
}

func UpdateUserData(pd *database.PlayerData) {
	userMu.Lock()
	defer userMu.Unlock()

	if users[pd.UUID] == nil {
		return
	}
	users[pd.UUID].Data = pd
}

func ResetUser(pl *player.Player) {
	ut := GetUser(pl)
	isBot := slices.Contains([]string{"Mark", "Sam", "Steven"}, pl.Name())
	pd := &database.PlayerData{
		UUID:                  pl.UUID(),
		UserId:                ut.Data.UserId,
		AlternativeMCAccounts: ut.Data.AlternativeMCAccounts,
		DeviceID:              ut.Data.DeviceID,
		HashedIP:              ut.Data.HashedIP,
		IPStoredSince:         ut.Data.IPStoredSince,
		Username:              pl.Name(),
		FirstJoin:             ut.Data.FirstJoin,
		LastJoin:              ut.Data.LastJoin,
		DeviceOS:              ut.Data.DeviceOS,
		ProtocolId:            ut.Data.ProtocolId,
		Statistics: database.Statistics{
			RankId: ut.Data.Statistics.RankId,
			Level:  1,
		},
		Cosmetics: database.Cosmetics{
			SelectedWoodType: block.OakWood(),
		},
		Settings: ut.Data.Settings,
	}

	if isBot {
		pd.DeviceOS = protocol.DeviceAndroid
	} else {
		pd.DeviceOS = utils.Session(pl).ClientData().DeviceOS
		pd.ProtocolId = utils.Session(pl).ClientData().GameVersion
	}

	u := &User{
		pl: pl,
		h:  pl.H(),

		cooldownMap: cooldown.NewMappedCoolDown[PlayerCoolDowns](),

		Data: pd,
	}
	users[pl.UUID()] = u

	_ = server.Database.SavePlayer(pd)
}

func (u *User) Player() *player.Player {
	return u.pl
}

func (u *User) H() *world.EntityHandle {
	return u.h
}

// IsCooldownActive sets a cooldown for the given type if it doesn't already exist or renews it if specified.
// It returns true if the cooldown is initially active (to be used for conditional command execution).
func (u *User) IsCooldownActive(cooldownType PlayerCoolDowns, duration time.Duration, renew, createIfInactive, sendMessage bool) bool {
	coolDown := u.cooldownMap[cooldownType]

	if coolDown == nil {
		coolDown = cooldown.NewCoolDown()
		u.cooldownMap[cooldownType] = coolDown
	}
	exists := coolDown.Active()
	if renew || createIfInactive && !coolDown.Active() {
		coolDown.Set(duration)
	}

	if sendMessage && exists {
		u.pl.Message(text.Colourf(language.Translate(u.pl).Commands.Error.CoolDown, coolDown.Remaining().Seconds()))
	}

	return exists
}

func (u *User) CoolDownTimeRemaining(cooldownType PlayerCoolDowns) time.Duration {
	coolDown := u.cooldownMap[cooldownType]
	if coolDown == nil || !coolDown.Active() {
		return 0
	}

	return coolDown.Remaining()
}

func ActiveBan(pd *database.PlayerData) *database.PunishmentData {
	checkBan := func(data *database.PlayerData) *database.PunishmentData {
		for _, ban := range slices.Concat(data.Punishments.Bans) {
			if ban.RemovedBy == "" && (ban.Permanent || ban.EndsAt.After(time.Now())) {
				return ban
			}
		}
		return nil
	}

	for _, alt := range pd.AlternativeMCAccounts {
		if ok := checkBan(utils.Panics(server.Database.FindPlayer(alt.UUID))); ok != nil {
			return ok
		}
	}

	return checkBan(pd)
}

func ActiveMute(pd *database.PlayerData) *database.PunishmentData {
	checkMute := func(data *database.PlayerData) *database.PunishmentData {
		for _, mute := range slices.Concat(data.Punishments.Mutes) {
			if mute.RemovedBy == "" && (mute.Permanent || mute.EndsAt.After(time.Now())) {
				return mute
			}
		}
		return nil
	}

	for _, alt := range pd.AlternativeMCAccounts {
		if ok := checkMute(utils.Panics(server.Database.FindPlayer(alt.UUID))); ok != nil {
			return ok
		}
	}

	return checkMute(pd)
}

func (u *User) AddItemWithHBConfig(preferredSlot int, it item.Stack) (n int, err error) {
	var category database.HotBarCategory

	// 1. Categorize the item
	switch it.Item().(type) {
	case item.Sword, item.Stick:
		category = database.Melee
	case block.Ladder:
		category = database.Ladder
	case item.GoldenApple, item.Snowball, block.TNT, item.EnderPearl, item.Bucket, item.Egg, block.Sponge, item.Compass:
		category = database.Utility
	case world.Block:
		category = database.Blocks
	case item.Bow:
		category = database.Bows
	case item.Potion:
		category = database.Potions
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

	// Try category slot
	if slot != -1 {
		sa, sb := it.AddStack(utils.Panics(u.pl.Inventory().Item(slot)))

		if sb.Count() != 0 {
			// Merge failed, check if different item
			if sa.Item() != sb.Item() {
				_ = u.pl.Inventory().SetItem(slot, sa)
				_, _ = u.pl.Inventory().AddItem(sb)
				return it.Count(), nil
			}
			// Merge failed but same item or can't move → fallback
			slot = -1
		} else {
			return it.Count(), u.pl.Inventory().SetItem(slot, sa)
		}
	}

	// Try preferredSlot (fallback)
	if slot == -1 && preferredSlot != -1 {
		sa, sb := it.AddStack(utils.Panics(u.pl.Inventory().Item(preferredSlot)))
		if sb.Count() == 0 {
			return it.Count(), u.pl.Inventory().SetItem(preferredSlot, sa)
		}
	}

	return u.pl.Inventory().AddItem(it)
}

func (u *User) SendScoreboard(numOfSpaces int) {
	if utils.Session(u.pl) == session.Nop {
		return
	}

	u.Scoreboard.Set(0, text.Colourf("%v%v", strings.Repeat(" ", numOfSpaces), font.Transform("SEASON 1")))
	u.pl.RemoveScoreboard()
	u.pl.SendScoreboard(u.Scoreboard)
}

func (u *User) DropItem(it item.Stack, tx *world.Tx) {
	ent := entity.NewItem(world.EntitySpawnOpts{
		Position: u.pl.Position(),
		Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1},
	}, it)
	tx.AddEntity(ent)
}

func (u *User) PlaySound(sound string, volume, pitch float64) {
	pos := u.pl.Position()
	pk := &packet.PlaySound{
		SoundName: sound,
		Position:  mgl32.Vec3{float32(pos.X()), float32(pos.Y()), float32(pos.Z())},
		Volume:    float32(volume),
		Pitch:     float32(pitch),
	}
	utils.WritePacket(utils.Session(u.pl), pk)
}

func (u *User) RefreshCape() {
	skin := u.pl.Skin()
	if cape, ok := database.GetCapeByIdentifier(u.Data.Cosmetics.SelectedCape); ok {
		skin.Cape = skin2.NewCape(64, 32)
		skin.Cape.Pix = database.CapeAsBytes(cape)
	} else {
		skin.Cape = skin2.NewCape(64, 32)
	}

	go func() {
		u.h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			e.(*player.Player).SetSkin(skin)
		})
	}()
}

type GameRuntimeData struct {
	BedWarsInfo  BedWarsInfo
	BuildFFAInfo BuildFFAInfo
}

type BedWarsInfo struct {
	Kills      int
	FinalKills int
	BedsBroken int
}

type BuildFFAInfo struct {
	Kills int
}

func (d GameRuntimeData) TotalBWKills() int {
	return d.BedWarsInfo.Kills + d.BedWarsInfo.FinalKills
}

type OldArmour struct {
	Helmet, ChestPlate, Leggings, Boots item.Stack
}
