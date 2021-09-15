package commands

const (
	CONF = "conf"
	SP   = "sp" // status power
	CP   = "cp" // control power
	ST   = "st" // status temperature
	SL   = "sl" // stauts level
	CL   = "cl" // control level
	CHS  = "chs"
	SHS  = "shs"
	CI   = "ci" // control identify
	SM   = "sm" //status motion
	SH   = "sh" //status humidity
)

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
