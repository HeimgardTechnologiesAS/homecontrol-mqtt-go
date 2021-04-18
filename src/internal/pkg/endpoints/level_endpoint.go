package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type LevelEndpoint struct {
	*endpoint
}

func NewLevelEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *LevelEndpoint {
	return &LevelEndpoint{
		endpoint: newEndpoint(
			"lev",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.CP: c.NewCommand(c.CP),
				c.SP: c.NewCommand(c.SP),
				c.CL: c.NewCommand(c.CL),
				c.SL: c.NewCommand(c.SL),
			}),
	}
}

func (obj *LevelEndpoint) SendStatus() {
	obj.SendFeedbackMessage(c.SP, obj.commands[c.SP].GetState())
	obj.SendFeedbackMessage(c.SL, obj.commands[c.SL].GetState())
}
