# github.com/HomeControlAS/homecontrol-mqtt-go

Alpha version

On Off endpoint example

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HomeControlAS/homecontrol-mqtt-go/commands"
	"github.com/HomeControlAS/homecontrol-mqtt-go/devices"
	"github.com/HomeControlAS/homecontrol-mqtt-go/endpoints"
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

func main() {

    // gw ip, mqtt device unique ID, gw username, gw pass, use secure TLS connection, custom name
	mqttDevice, err := devices.NewMqttDevice("192.168.8.1", "test_dev12345", "hc", "admin", true, "mqtt_device")
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

	quitCh := make(chan error)
	setSignalInterrupt(quitCh)
	err = mqttDevice.RunForever(quitCh)
	if err != nil {
		log.Printf("Device stopped unexpectedly. %s", err.Error())
	}
}

func setSignalInterrupt(quitCh chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quitCh <- nil
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}()
}
```
