package devices

import "homecontrol-mqtt-go/internal/pkg/endpoints"

type Device interface {
	Connect() error
	Disconnect()
	RunForever() error
	GetQuitCh() chan error
	AddEndpoint(enp endpoints.Endpoint)
	SendConfigs()
}
