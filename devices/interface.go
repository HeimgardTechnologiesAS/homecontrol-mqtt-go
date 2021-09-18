package devices

import "github.com/HomeControlAS/homecontrol-mqtt-go/endpoints"

// Device interface that collects methods needed to connect MQTT device to HC Gateway
type Device interface {
	// Connect connects Device object to HC gateway
	Connect() error
	// Disconnect disconnects Device from HC gateway
	Disconnect()
	// AddEndpoint adds new endpoint to MQTT Device
	AddEndpoint(enp endpoints.Endpoint)
	// GetEndpoint returns endpoint with given ID
	GetEndpoint(uid string) endpoints.Endpoint
	// RegisterOnConnectCb registers handler that is invoked when the connection is established
	RegisterOnConnectCb(cb func())
	// RegisterOnConnectionLostCb registers handler that is invoked when connection is lost
	RegisterOnConnectionLostCb(cb func(err error))
}
