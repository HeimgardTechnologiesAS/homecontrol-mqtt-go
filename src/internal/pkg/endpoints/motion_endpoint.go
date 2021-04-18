package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type MotionEndpoint struct {
	*endpoint
}

func NewMotionEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *MotionEndpoint {
	return &MotionEndpoint{
		endpoint: newEndpoint(
			"mot",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.SM: c.NewCommand(c.SM),
			}),
	}
}

func (obj *MotionEndpoint) SendStatus() {
	obj.SendFeedbackMessage(c.SM, obj.commands[c.SM].GetState())
}
