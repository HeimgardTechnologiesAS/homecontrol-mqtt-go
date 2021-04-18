package endpoints

import (
	c "homecontrol-mqtt-go/internal/pkg/commands"
)

type IdentifyEndpoint struct {
	*endpoint
}

func NewIdentifyEndpoint(
	epId string,
	epName string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *IdentifyEndpoint {
	return &IdentifyEndpoint{
		endpoint: newEndpoint(
			"id",
			"60",
			epId,
			epName,
			onStateChange,
			map[string]c.Command{
				c.CI: c.NewCommand(c.CI),
			}),
	}
}

func (obj *IdentifyEndpoint) SendStatus() {
	// do nothing
}
