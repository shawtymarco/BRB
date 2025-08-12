package bedwars

import (
	"server/server/game"
	"server/server/living"
	"server/server/utils"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type GeneratorSettings struct {
	Active bool
	Game   *BedWars

	Resource Resource
	Name     string
	Tier     int

	Cap       int
	SpawnRate time.Duration
}

func (gs GeneratorSettings) New(pos mgl64.Vec3, tx *world.Tx) *GeneratorBlockType {
	t := &GeneratorBlockType{GeneratorSettings: gs}
	if gs.Game != nil {
		t.Team = gs.Game.NearestTeam(pos)
	}

	conf := living.Config{
		EntityType: t,
		MaxHealth:  1,
		Speed:      0.1,
		MovementComputer: &entity.MovementComputer{
			Gravity:           0,
			Drag:              0,
			DragBeforeGravity: false,
		},
		Drops: []living.Drop{living.NewDropWithStack(item.NewStack(gs.Resource.Item(), 1))},
	}

	gb := tx.AddEntity(world.EntitySpawnOpts{Position: pos.Add(mgl64.Vec3{0, 2, 0})}.New(conf.EntityType, conf)).(*GeneratorBlockType)
	gb.SetImmobile(true, tx)

	_, ok1 := gb.Resource.Block().(block.Air)
	_, ok2 := gb.Resource.Block().(block.Diamond)
	gb.WithVariant(int32(lo.If(ok1 || gs.Game != nil, 0).ElseIf(ok2, 1).Else(2)))
	gb.lastSpawn = time.Now()
	return gb
}

type GeneratorBlockType struct {
	living.NopLivingType

	*living.Living
	GeneratorSettings
	Team *game.Team

	tick      time.Duration
	lastSpawn time.Time

	queue []*player.Player
}

func (*GeneratorBlockType) EncodeEntity() string   { return "bedwars:generator_block" }
func (*GeneratorBlockType) NetworkOffset() float64 { return 0.49 }
func (*GeneratorBlockType) BBox(world.Entity) cube.BBox {
	return utils.ZeroBox
}

func (b *GeneratorBlockType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	b.Living = b.NopLivingType.Open(tx, handle, data).(*living.Living)
	return b
}

func (b *GeneratorBlockType) PlayersWithin(tx *world.Tx) []*player.Player {
	var res []*player.Player
	for e := range tx.Entities() {
		if pl, ok := e.(*player.Player); ok && pl.GameMode() != world.GameModeSpectator && utils.Distance(b.Position(), e.Position()) <= 4 {
			res = append(res, pl)
		}
	}
	return res
}

func (b *GeneratorBlockType) CountResourcesWithin() (res int) {
	for e := range b.Tx().Entities() {
		if ent, ok := e.(*entity.Ent); ok && utils.Distance(b.Position(), e.Position()) <= 4 && e.H().Type() == entity.ItemType {
			if beh, ok := ent.Behaviour().(*entity.ItemBehaviour); ok && beh.Item().Item() == b.Resource.Item() {
				res += beh.Item().Count()
			}
		}
	}
	return res
}

func (b *GeneratorBlockType) ResourcesWithin(tx *world.Tx) []*entity.Ent {
	var res []*entity.Ent

	for e := range tx.Entities() {
		if ent, ok := e.(*entity.Ent); ok && utils.Distance(b.Position(), e.Position()) <= 4 && e.H().Type() == entity.ItemType {
			if beh, ok := ent.Behaviour().(*entity.ItemBehaviour); ok && beh.Item().Item() == b.Resource.Item() {
				res = append(res, ent)
			}
		}
	}
	return res
}

func (b *GeneratorBlockType) UpdateQueue(tx *world.Tx) {
	playersWithin := b.PlayersWithin(tx)

	for _, p := range playersWithin {
		if !slices.ContainsFunc(b.queue, func(p2 *player.Player) bool {
			return p2.UUID() == p.UUID()
		}) {
			b.AddPlayer(p)
		}
	}

	for _, p := range b.queue {
		if !slices.ContainsFunc(playersWithin, func(p2 *player.Player) bool {
			return p2.UUID() == p.UUID()
		}) {
			b.RemovePlayer(p)
		}
	}
}

func (b *GeneratorBlockType) Next() *player.Player {
	pl := b.queue[0]
	b.queue = b.queue[1:]
	b.queue = append(b.queue, pl)
	return pl
}

func (b *GeneratorBlockType) AddPlayer(pl *player.Player) {
	b.queue = append(b.queue, pl)
}

func (b *GeneratorBlockType) RemovePlayer(pl *player.Player) {
	for i, p := range b.queue {
		if p == pl {
			b.queue = append(b.queue[:i], b.queue[i+1:]...)
			break
		}
	}
}

func (b *GeneratorBlockType) Hurt(dmg float64, src world.DamageSource) (float64, bool) {
	return 0, false
}

func (b *GeneratorBlockType) Tick(tx *world.Tx, current int64) {
	if !b.Active {
		b.Living.Tick(tx, current+1)
		return
	}

	if b.Team != nil {
		b.Tier = b.Team.Upgrades.GeneratorTier + 1
		b.updateTier()
	}

	remainingDur := lo.If(b.CountResourcesWithin() >= b.Cap, b.SpawnRate).Else(b.SpawnRate - time.Now().Sub(b.lastSpawn))
	if _, ok := b.Resource.Block().(block.Air); !ok && b.Team == nil {
		b.SetNameTag(text.Colourf("<bold><yellow>Tier <red>%v</red></yellow></bold>\n%v\n<yellow>Spawns in <red>%.1f</red> seconds</yellow>", utils.IntToRoman(b.Tier), b.Name, remainingDur.Seconds()), tx)
	}

	if b.CountResourcesWithin() < b.Cap {
		if remainingDur <= 0 {
			b.DropItems(tx)
			b.lastSpawn = time.Now()
		}
		b.tick += 50 * time.Millisecond
	}
	b.Living.Tick(tx, current+1)
}

func (b *GeneratorBlockType) updateTier() {
	genTier := b.Team.Upgrades.GeneratorTier
	switch b.Resource {
	case Iron:
		switch genTier {
		case 1:
			b.SpawnRate = 300 * time.Millisecond
			b.Cap = 64
		case 2:
			b.SpawnRate = 200 * time.Millisecond
		case 3:
		case 4:
			b.SpawnRate = 150 * time.Millisecond
		}
	case Gold:
		switch genTier {
		case 1:
			b.SpawnRate = 3000 * time.Millisecond
			b.Cap = 16
		case 2:
			b.SpawnRate = 2000 * time.Millisecond
		case 3:
		case 4:
			b.SpawnRate = 1350 * time.Millisecond
		}
	case Emerald:
		if genTier == 4 {
			b.SpawnRate = 15 * time.Second
		}
	default:
	}
}
