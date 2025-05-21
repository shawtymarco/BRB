package user

type PlayerCoolDowns int

const (
	INTERACT PlayerCoolDowns = iota
	CommandPing
	CommandHub
)
