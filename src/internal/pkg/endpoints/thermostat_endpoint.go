package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type ThermostatEndpoint struct {
	*endpoint
}

func NewThermostatEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *ThermostatEndpoint {
	return &ThermostatEndpoint{

		endpoint: newEndpoint(
			"thrmstt",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.ST:  c.NewCommand(c.ST),
				c.CHS: c.NewCommand(c.CHS),
				c.SHS: c.NewCommand(c.SHS),
			}),
	}
}

func (obj *ThermostatEndpoint) SendStatus() {
	obj.SendFeedbackMessage(c.ST, obj.commands[c.ST].GetState())
	obj.SendFeedbackMessage(c.SHS, obj.commands[c.SHS].GetState())
}
