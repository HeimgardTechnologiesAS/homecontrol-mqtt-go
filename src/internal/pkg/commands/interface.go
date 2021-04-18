package commands

type Command interface {
	SetState(currentState string)
	GetState() string
}
