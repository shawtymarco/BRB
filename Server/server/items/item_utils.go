package items

import (
	"image"
	"image/png"
	"os"
	"server/server/utils"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var ItemHandlers = map[int]ItemHandler{}

type ClickType uint8

const (
	OnStartBreak ClickType = iota
	OnPunchAir
	OnItemUse
	OnItemUseOnBlock
	OnItemUseOnEntity
)

type ItemHandler interface {
	Stack() item.Stack
	InteractClick(click ClickType, player *player.Player, pos *cube.Pos)
}

func Texture(file string) image.Image {
	texture := utils.Panics(os.OpenFile(file, os.O_RDONLY, os.ModePerm))
	defer utils.Panic(texture.Close())
	return utils.Panics(png.Decode(texture))
}

type ItemAction int

const (
	Wand ItemAction = iota
	LeaveGame
	TeamSelector
	Cosmetics

	Beam
	WhiteSword
	MiniSun
	SpawnMage
)

type NopItem struct {
}

func (NopItem) Name() string {
	return ""
}

func (NopItem) Texture() image.Image {
	return nil
}

func init() {
	item.RegisterEnchantment(1000, Glitter{})
}

type Glitter struct{}

func (Glitter) Name() string {
	return ""
}

func (Glitter) MaxLevel() int {
	return 1
}

func (Glitter) Cost(int) (int, int) {
	return 0, 1
}

func (Glitter) Rarity() item.EnchantmentRarity {
	return CustomRarityServerUse{}
}

func (Glitter) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (Glitter) CompatibleWithItem(_ world.Item) bool {
	return true
}

type CustomRarityServerUse struct{}

func (CustomRarityServerUse) Name() string { return "" }
func (CustomRarityServerUse) Cost() int    { return 1 }
func (CustomRarityServerUse) Weight() int  { return 0 }
