package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homecontrol-mqtt-go/internal/pkg/commands"
	"homecontrol-mqtt-go/internal/pkg/devices"
	"homecontrol-mqtt-go/internal/pkg/endpoints"
)

func ep1StateChange(ep endpoints.Endpoint, cmd string, msg string) {
	fmt.Printf("GOT STATE UPDATE EP %s, CMD %s, MSG %s\n", ep.GetID(), cmd, msg)
	if msg == "1" {
		ep.SendFeedbackMessage(commands.SP, "1")
	} else {
		ep.SendFeedbackMessage(commands.SP, "0")
	}
}

func main() {
	quit := make(chan int)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quit <- 0
		fmt.Println("Done")
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}()

	mqttDevice := devices.NewMqttDevice("192.168.8.1", 1883, "ttdev123", "hc", "admin", false)

	ep1 := endpoints.NewOnOffEndpoint("ep1", ep1StateChange)
	ep2 := endpoints.NewOnOffEndpoint("ep2", ep1StateChange)

	mqttDevice.AddEndpoint(ep1)
	mqttDevice.AddEndpoint(ep2)

	mqttDevice.Connect()
	defer mqttDevice.Disconnect()

	mqttDevice.RunForever(quit)
}
