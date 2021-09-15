package endpoints

import (
	"fmt"
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type LevelEndpoint struct {
	*endpoint
}

func NewLevelEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
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

func (obj *LevelEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.SP, obj.commands[c.SP].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SP feedback: %s", obj.GetID(), err)
	}
	err = obj.SendFeedbackMessage(c.SL, obj.commands[c.SL].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SL feedback: %s", obj.GetID(), err)
	}
	return nil
}
