package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homecontrol-mqtt-go/internal/pkg/commands"
	"homecontrol-mqtt-go/internal/pkg/devices"
	"homecontrol-mqtt-go/internal/pkg/endpoints"
)

func ep1StateChange(ep endpoints.Endpoint, cmd string, msg string) {
	log.Printf("Message received. Endpoint: %s, Command: %s, Message: %s\n", ep.GetID(), cmd, msg)
	if msg == "1" {
		_ = ep.SendFeedbackMessage(commands.SP, "1")
	} else {
		_ = ep.SendFeedbackMessage(commands.SP, "0")
	}
}

func main() {

	mqttDevice, err := devices.NewMqttDevice("192.168.8.1", "ttdev1234", "hc", "admin", true)
	if err != nil {
		log.Printf("failed to create MQTT device: %s\n", err.Error())
		return
	}

	ep1 := endpoints.NewOnOffEndpoint("ep1", ep1StateChange)

	mqttDevice.AddEndpoint(ep1)

	err = mqttDevice.Connect()
	if err != nil {
		log.Printf("failed to connect %s\n", err.Error())
		return
	}
	defer mqttDevice.Disconnect()

	setSignalInterrupt(mqttDevice.GetQuitCh())
	err = mqttDevice.RunForever()
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
