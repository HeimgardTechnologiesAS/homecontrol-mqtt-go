package commands

type command struct {
	cmdType string
	state   string
}

func NewCommand(cmdType string) Command {
	return &command{
		cmdType: cmdType,
		state:   "0",
	}
}

func (obj *command) SetState(currentState string) {
	obj.state = currentState
}

func (obj *command) GetState() string {
	return obj.state
}
