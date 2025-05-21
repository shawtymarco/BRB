package user

type SpectatorGamemode struct{}

func (sg SpectatorGamemode) AllowsEditing() bool {
	return false
}

func (sg SpectatorGamemode) AllowsTakingDamage() bool {
	return false
}

func (sg SpectatorGamemode) CreativeInventory() bool {
	return false
}

func (sg SpectatorGamemode) HasCollision() bool {
	return true
}

func (sg SpectatorGamemode) AllowsFlying() bool {
	return true
}

func (sg SpectatorGamemode) AllowsInteraction() bool {
	return false
}

func (sg SpectatorGamemode) Visible() bool {
	return false
}
