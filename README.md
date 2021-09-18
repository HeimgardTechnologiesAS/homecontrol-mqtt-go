# github.com/HomeControlAS/homecontrol-mqtt-go

Alpha version

On Off endpoint example

```go
package main

import (
	"log"

	"github.com/HomeControlAS/homecontrol-mqtt-go/commands"
	"github.com/HomeControlAS/homecontrol-mqtt-go/devices"
	"github.com/HomeControlAS/homecontrol-mqtt-go/endpoints"
	"github.com/HomeControlAS/homecontrol-mqtt-go/utils"
)

func ep1StateChange(ep endpoints.Endpoint, cmd string, msg string, err error) {
	log.Printf("Message received. Endpoint: %s, Command: %s, Message: %s\n", ep.GetID(), cmd, msg)
	if err != nil {
		log.Printf("error when handling state change: %s", err)
		return
	}
	// msg can be equal to "0" or "1", just send it back, since there is nothing to control
	err = ep.SendFeedbackMessage(commands.SP, msg)
	if err != nil {
		log.Printf("error while sending feedback message %s", err)
	}
}

func onConnectionLostEvent(err error) {
	log.Printf("Event, Connection lost %s", err)
	utils.TriggerSignalInterrupt()
}

func main() {

	mqttDevice, err := devices.NewMqttDevice("192.168.8.1", "newDev1", "hc", "admin", true, "mqtt_device")
	mqttDevice.RegisterOnConnectionLostCb(onConnectionLostEvent)
	if err != nil {
		log.Printf("failed to create MQTT device: %s\n", err.Error())
		return
	}

	mqttDevice.AddEndpoint(endpoints.NewOnOffEndpoint("ep1", "On_Off", ep1StateChange))

	err = mqttDevice.Connect()
	if err != nil {
		log.Printf("failed to connect %s\n", err.Error())
		return
	}
	defer mqttDevice.Disconnect()
	utils.WaitForSignalInterrupt()
}
```
