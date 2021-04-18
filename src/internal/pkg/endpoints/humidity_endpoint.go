package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type HumidityEndpoint struct {
	*endpoint
}

func NewHumidityEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *HumidityEndpoint {
	return &HumidityEndpoint{
		endpoint: newEndpoint(
			"hum",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.SH: c.NewCommand(c.SH),
			}),
	}
}

func (obj *HumidityEndpoint) SendStatus() {
	obj.SendFeedbackMessage(c.SH, obj.commands[c.SH].GetState())
}
