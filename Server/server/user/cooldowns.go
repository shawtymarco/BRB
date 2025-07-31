package user

type PlayerCoolDowns int

const (
	Interact PlayerCoolDowns = iota
	CriticalExtraHit
	Switching
	Chat
	CommandPing
	CommandHub
	CommandLink
)
