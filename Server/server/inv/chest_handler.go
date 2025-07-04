package inv

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

type ChestUIHandler struct {
	inventory.NopHandler
	Inventory *inventory.Inventory
	Funcs     []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory)
}

func (h ChestUIHandler) HandleTake(ctx *event.Context[inventory.Holder], slot int, stack item.Stack) {
	ctx.Val()
	if h.Funcs[0] != nil {
		h.Funcs[0](ctx, slot, stack, h.Inventory)
	}
}

func (h ChestUIHandler) HandlePlace(ctx *event.Context[inventory.Holder], slot int, stack item.Stack) {
	if h.Funcs[1] != nil {
		h.Funcs[1](ctx, slot, stack, h.Inventory)
	}
}

func (h ChestUIHandler) HandleDrop(ctx *event.Context[inventory.Holder], slot int, stack item.Stack) {
	if h.Funcs[2] != nil {
		h.Funcs[2](ctx, slot, stack, h.Inventory)
	}
}
