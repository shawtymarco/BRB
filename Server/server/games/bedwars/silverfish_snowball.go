package bedwars

import (
	"server/server/itemutil"
	user2 "server/server/user"

	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.SilverfishSnowball, SilverfishSnowballItem{})
}

func newSilverfishSnowball(opts world.EntitySpawnOpts, game *BedWars, owner world.Entity) *world.EntityHandle {
	conf := snowballConf
	conf.Owner = owner.H()
	conf.Hit = func(e *entity.Ent, tx *world.Tx, target trace.Result) {
		if ownerEntity, ok := owner.H().Entity(tx); ok {
			if pl, ok := ownerEntity.(*player.Player); ok {
				NewSilverfish(target.Position(), game, game.PlayerTeam(pl), owner.H(), tx)
			}
		}
	}
	return opts.New(entity.SnowballType, conf)
}

var snowballConf = entity.ProjectileBehaviourConfig{
	Gravity:       0.03,
	Drag:          0.01,
	Particle:      particle.SnowballPoof{},
	ParticleCount: 6,
}

type SilverfishSnowballItem struct {
	item.Snowball
}

func NewSilverfishSnowballItem() item.Stack {
	return item.NewStack(SilverfishSnowballItem{}, 1).WithValue("special_item", int16(itemutil.SilverfishSnowball))
}

func (s SilverfishSnowballItem) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pl := user.(*player.Player)
	u := user2.GetUser(pl)
	g := Games[u.Game.ID()]

	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(1.5)}
	tx.AddEntity(newSilverfishSnowball(opts, g, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}
