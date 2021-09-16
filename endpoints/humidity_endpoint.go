package endpoints

import (
	"fmt"

	c "github.com/HomeControlAS/homecontrol-mqtt-go/commands"
)

type HumidityEndpoint struct {
	*endpoint
}

func NewHumidityEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
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

func (obj *HumidityEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.SH, obj.commands[c.SH].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SH feedback: %s", obj.GetID(), err)
	}
	return nil
}
