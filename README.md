# homecontrol-mqtt-go

### Setup
use `go get github.com/HomeControlAS/homecontrol-mqtt-go`

or 

* clone the git repo
* open directory in vscode and docker devcontainer
* go to examples
* modify GW IP and MQTT credentials
* in the `main.go` directory run `go run .`


Currently in beta test phase

On/Off device example 

```go
package main

import (
	"log"

	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/commands"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/devices"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/endpoints"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/utils"
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

