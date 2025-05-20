package living

import (
	"fmt"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type NopLivingType struct{}

func (n NopLivingType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Living{
		livingData: data.Data.(*livingData),
		tx:         tx,
		handle:     handle,
		data:       data,
	}
}

func (NopLivingType) EncodeEntity() string {
	return utils.EnumPanic[string](fmt.Errorf("implement me"))
}

func (NopLivingType) BBox(_ world.Entity) cube.BBox {
	return cube.BBox{}
}

func (NopLivingType) DecodeNBT(m map[string]any, data *world.EntityData) {
	data.Data = m
}

func (NopLivingType) EncodeNBT(data *world.EntityData) map[string]any {
	return map[string]any{
		"data": data,
	}
}
