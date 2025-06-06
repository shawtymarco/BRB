package game

type Stage int

const (
	Waiting Stage = iota
	Starting
	Running
	Ending
	Terminated
)
