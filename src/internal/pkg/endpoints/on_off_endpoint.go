package endpoints

import (
	"fmt"
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type OnOffEndpoint struct {
	*endpoint
}

func NewOnOffEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
) *OnOffEndpoint {
	return &OnOffEndpoint{
		endpoint: newEndpoint(
			"pwr",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.CP: c.NewCommand(c.CP),
				c.SP: c.NewCommand(c.SP),
			}),
	}
}

func (obj *OnOffEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.SP, obj.commands[c.SP].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SP feedback: %s", obj.GetID(), err)
	}
	return nil
}
