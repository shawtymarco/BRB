package bedwars

import (
	"fmt"
	"math"
	"server/server/game"
	"server/server/living"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Silverfish struct {
	living.NopLivingType

	*living.Living

	tick time.Duration

	game *BedWars
	team *game.Team

	owner  *world.EntityHandle
	target *world.EntityHandle

	lastAttack time.Time
}

func NewSilverfish(pos mgl64.Vec3, game *BedWars, team *game.Team, owner *world.EntityHandle, tx *world.Tx) *Silverfish {
	t := &Silverfish{
		game:  game,
		team:  team,
		owner: owner,
	}

	conf := living.Config{
		EntityType: t,
		MaxHealth:  5,
		Speed:      0.3,
		MovementComputer: &entity.MovementComputer{
			Gravity:           0.5,
			Drag:              0,
			DragBeforeGravity: false,
		},
		Drops: []living.Drop{},
	}

	s := tx.AddEntity(world.EntitySpawnOpts{Position: pos}.New(conf.EntityType, conf)).(*Silverfish)

	ownerPl, _ := owner.Entity(tx)
	s.SetNameTag(text.Colourf(fmt.Sprintf("<%v>[%v]</%v> %v's Silverfish", team.Colour(), strings.ToUpper(team.Colour()), team.Colour(), ownerPl.(*player.Player).Name())), tx)

	return s
}

func (*Silverfish) EncodeEntity() string { return "minecraft:silverfish" }

func (*Silverfish) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.2, 0, -0.25, 0.2, 0.3, 0.25)
}
func (s *Silverfish) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	s.Living = s.NopLivingType.Open(tx, handle, data).(*living.Living)
	return s
}

func (s *Silverfish) Hurt(dmg float64, src world.DamageSource) (float64, bool) {
	d, i := s.Living.Hurt(dmg, src)

	if src2, ok := src.(entity.AttackDamageSource); ok {
		if pl, ok := src2.Attacker.(*player.Player); ok && s.game.PlayerTeam(pl) == s.team {
			return 0, false
		}
	}

	return d, i
}

func (s *Silverfish) Tick(tx *world.Tx, current int64) {
	if s.target == nil {
		if nearestEnemy := s.findNearestEnemy(tx); nearestEnemy != nil {
			s.target = nearestEnemy.H()
		}
		s.Living.Tick(tx, current+1)
		return
	}
	te, ok := s.target.Entity(tx)
	if !ok || te.(entity.Living).Dead() {
		s.target = s.findNearestEnemy(tx).H()
		s.Living.Tick(tx, current+1)
		return
	}

	if s.Position().Y() < float64(s.game.MapConfig().Void) {
		s.Kill(entity.VoidDamageSource{})
		return
	}

	s.LookAt(te.Position(), tx)

	if utils.Distance(s.Position(), te.Position()) > 1 {
		s.MoveToTarget(te.Position(), 1, tx)
	} else if time.Now().Sub(s.lastAttack) > 1*time.Second {
		s.lastAttack = time.Now()
		if ownerEnt, ok := s.owner.Entity(tx); ok {
			te.(entity.Living).Hurt(5, entity.AttackDamageSource{Attacker: ownerEnt})
		}
	}

	s.tick += 50 * time.Millisecond
	s.Living.Tick(tx, current+1)
}

func (s *Silverfish) findNearestEnemy(tx *world.Tx) (nearestEnemy *player.Player) {
	var minDist = math.MaxFloat64

	for _, e := range s.game.ActivePlayers() {
		if ent, ok := e.Entity(tx); ok {
			pl, _ := ent.(*player.Player)
			if pl.GameMode() != world.GameModeSpectator && s.game.PlayerTeam(pl) != s.team {
				dist := utils.Distance(s.Position(), pl.Position())
				if dist < minDist {
					minDist = dist
					nearestEnemy = pl
				}
			}
		}
	}

	return nearestEnemy
}
