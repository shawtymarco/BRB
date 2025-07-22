package user

type PlayerCoolDowns int

const (
	Interact PlayerCoolDowns = iota
	Switching
	Chat
	CommandPing
	CommandHub
	CommandLink
)
