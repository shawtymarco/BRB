package listener

import (
	"fmt"
	"reflect"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type WorldHandler struct {
	world.NopHandler
}

func (h WorldHandler) HandleSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3) {
	fmt.Println(reflect.TypeOf(s))
	if _, ok := s.(sound.Attack); ok {
		ctx.Cancel()
	}
}
