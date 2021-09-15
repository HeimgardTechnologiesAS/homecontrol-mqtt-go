package endpoints

import (
	"fmt"
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type MotionEndpoint struct {
	*endpoint
}

func NewMotionEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string, err error),
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

func (obj *MotionEndpoint) SendStatus() error {
	err := obj.SendFeedbackMessage(c.SM, obj.commands[c.SM].GetState())
	if err != nil {
		return fmt.Errorf("endpoint [%s], failed to send SM feedback: %s", obj.GetID(), err)
	}
	return nil
}
