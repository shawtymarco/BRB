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

type TWGamemode struct{}

func (twg TWGamemode) AllowsEditing() bool {
	return false
}

func (twg TWGamemode) AllowsTakingDamage() bool {
	return true
}

func (twg TWGamemode) CreativeInventory() bool {
	return false
}

func (twg TWGamemode) HasCollision() bool {
	return true
}

func (twg TWGamemode) AllowsFlying() bool {
	return false
}

func (twg TWGamemode) AllowsInteraction() bool {
	return true
}

func (twg TWGamemode) Visible() bool {
	return true
}
