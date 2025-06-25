package user

type PlayerCoolDowns int

const (
	Interact PlayerCoolDowns = iota
	Chat
	CommandPing
	CommandHub
	CommandLink
)
