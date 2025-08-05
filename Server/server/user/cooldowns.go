package user

type PlayerCoolDowns int

const (
	Interact PlayerCoolDowns = iota
	NoPVP
	Chat
	CommandPing
	CommandHub
	CommandLink
)
