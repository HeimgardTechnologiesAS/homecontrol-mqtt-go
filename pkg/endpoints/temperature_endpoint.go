package endpoints

import (
	"fmt"

	c "github.com/HomeControlAS/homecontrol-mqtt-go/pkg/commands"
)

type TemperatureEndpoint struct {
	*endpoint
}

func NewTemperatureEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
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

func (obj *TemperatureEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.ST, obj.commands[c.ST].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send ST feedback: %s", obj.GetID(), err)
	}
	return nil
}
