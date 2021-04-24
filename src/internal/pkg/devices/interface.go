package devices

import "homecontrol-mqtt-go/internal/pkg/endpoints"

// Device interface that collects methods needed to connect MQTT device to HC Gateway
type Device interface {
	// Connect connects Device object to HC gateway
	Connect() error
	// Disconnect disconnects Device from HC gateway
	Disconnect()
	// RunForever runs infinite loop if MQTT Device should listen forever
	RunForever() error
	// GetQuitCh returns Quit Channel needed to stop RunForever loop
	GetQuitCh() chan error
	// AddEndpoint adds new endpoint to MQTT Device
	AddEndpoint(enp endpoints.Endpoint)
}
