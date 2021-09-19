package main

import (
	"log"
	"strconv"
	"time"

	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/commands"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/devices"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/endpoints"
	"github.com/HomeControlAS/homecontrol-mqtt-go/pkg/utils"
)

var fakeTemperature = 0

func onConnectEvent() {
	log.Printf("Event, Connected")
}

func onConnectionLostEvent(err error) {
	log.Printf("Event, Connection lost %s", err)
	utils.TriggerSignalInterrupt()
}

func main() {

	mqttDevice, err := devices.NewMqttDevice("192.168.8.1", "temp_dev", "hc", "admin", true, "mqtt_temp_device")
	mqttDevice.RegisterOnConnectCb(onConnectEvent)
	mqttDevice.RegisterOnConnectionLostCb(onConnectionLostEvent)
	if err != nil {
		log.Printf("failed to create MQTT device: %s\n", err.Error())
		return
	}

	// this endpoint supports only status temperature command -> onStateChanged handler is not needed
	mqttDevice.AddEndpoint(endpoints.NewTemperatureEndpoint("ep", "tmp", nil))

	err = mqttDevice.Connect()
	if err != nil {
		log.Printf("failed to connect %s\n", err.Error())
		return
	}
	defer mqttDevice.Disconnect()

	e := mqttDevice.GetEndpoint("ep")

	if e != nil {
		log.Printf("runing fake temperature reporter")
		// fake temperature reporter
		go func(e endpoints.Endpoint, cnt int) {
			for {
				err := e.SendFeedbackMessage(commands.ST, strconv.Itoa(fakeTemperature))
				if err != nil {
					log.Printf("failed to update temperature")
				}
				// increase temperature every 5 seconds and send
				fakeTemperature++
				time.Sleep(5 * time.Second)
			}
		}(e, fakeTemperature)
	}

	// exit on ctrl+c
	utils.WaitForSignalInterrupt()
}
