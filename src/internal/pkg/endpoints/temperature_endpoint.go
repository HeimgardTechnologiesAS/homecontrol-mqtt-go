package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type TemperatureEndpoint struct {
	*endpoint
}

func NewTemperatureEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *TemperatureEndpoint {
	return &TemperatureEndpoint{
		endpoint: newEndpoint(
			"tmp",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.ST: c.NewCommand(c.ST),
			}),
	}
}

func (obj *TemperatureEndpoint) SendStatus() {
	obj.SendFeedbackMessage(c.ST, obj.commands[c.ST].GetState())
}
