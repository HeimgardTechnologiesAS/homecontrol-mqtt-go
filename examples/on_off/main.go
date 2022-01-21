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

	err = ep.SendFeedbackMessage(commands.SP, msg)
	if err != nil {
		log.Printf("error while sending feedback message %s", err)
	}
}

func ep2StateChange(ep endpoints.Endpoint, cmd string, msg string, err error) {
	log.Printf("Message received. Endpoint: %s, Command: %s, Message: %s\n", ep.GetID(), cmd, msg)
	if err != nil {
		log.Printf("error when handling state change: %s", err)
		return
	}
	//msg can be equal to "0" or "1", just send it back, since there is nothing to control
	if cmd == commands.CP {
		err = ep.SendFeedbackMessage(commands.SP, msg)
		if err != nil {
			log.Printf("error while sending feedback message %s", err)
		}
	}
	if cmd == commands.CL {
		err = ep.SendFeedbackMessage(commands.SL, msg)
		if err != nil {
			log.Printf("error while sending feedback message %s", err)
		}
	}
}

func onConnectionLostEvent(err error) {
	log.Printf("Event, Connection lost %s", err)
	utils.TriggerSignalInterrupt()
}

func main() {

	mqttDevice, err := devices.NewMqttDevice("ENTER_IP", "test_dev", "MQTT_USERNAME", "MQTT_PASS", true, "mqtt_device")
	if err != nil {
		log.Printf("failed to create MQTT device: %s\n", err.Error())
		return
	}
	mqttDevice.RegisterOnConnectionLostCb(onConnectionLostEvent)

	mqttDevice.AddEndpoint(endpoints.NewOnOffEndpoint("ep1", "SmartPlug", ep1StateChange))
	mqttDevice.AddEndpoint(endpoints.NewLevelEndpoint("ep2", "Bulb", ep2StateChange))

	err = mqttDevice.Connect()
	if err != nil {
		log.Printf("failed to connect %s\n", err.Error())
		return
	}
	defer mqttDevice.Disconnect()

	utils.WaitForSignalInterrupt()
}
