package endpoints

import (
	"fmt"

	c "github.com/HomeControlAS/homecontrol-mqtt-go/commands"
)

type ThermostatEndpoint struct {
	*endpoint
}

func NewThermostatEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
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

func (obj *ThermostatEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.ST, obj.commands[c.ST].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send ST feedback: %s", obj.GetID(), err)
	}
	err = obj.SendFeedbackMessage(c.SHS, obj.commands[c.SHS].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SHS feedback: %s", obj.GetID(), err)
	}
	return nil
}
