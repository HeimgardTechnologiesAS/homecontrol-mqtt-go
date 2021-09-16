package devices

import "github.com/HomeControlAS/homecontrol-mqtt-go/endpoints"

// Device interface that collects methods needed to connect MQTT device to HC Gateway
type Device interface {
	// Connect connects Device object to HC gateway
	Connect() error
	// Disconnect disconnects Device from HC gateway
	Disconnect()
	// RunForever runs infinite loop if MQTT Device should listen forever
	RunForever(quitC chan error) error
	// AddEndpoint adds new endpoint to MQTT Device
	AddEndpoint(enp endpoints.Endpoint)
}
