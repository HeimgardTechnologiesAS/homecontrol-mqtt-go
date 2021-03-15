package devices

import (
	"errors"
	"fmt"
	"homecontrol-mqtt-go/internal/pkg/endpoints"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttDevice struct {
	uid       string
	client    mqtt.Client
	epZero    *endpoints.ZeroEndpoint
	endpoints map[string]endpoints.Endpoint
}

func parseEndpointID(token string) string {
	fDelimeterIn := strings.Index(token, "/")
	if fDelimeterIn > -1 {
		enpIn := token[fDelimeterIn+1:]
		sDelimeterInd := strings.Index(enpIn, "/")
		if sDelimeterInd > 0 {
			return enpIn[0:sDelimeterInd]
		}
	}
	return ""
}

func parseCommand(token string) string {
	lDelimeterIn := strings.LastIndex(token, "/")
	if lDelimeterIn > -1 {
		return token[lDelimeterIn+1:]
	}
	return ""
}

func NewMqttDevice(hostname string, port int, uid string, username string, password string, isSecure bool) *MqttDevice {

	device := &MqttDevice{
		uid:       uid,
		endpoints: make(map[string]endpoints.Endpoint),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", hostname, port))
	opts.SetClientID(uid)
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(device.onMessageHandler)
	opts.OnConnect = device.onConnectHandler
	opts.OnConnectionLost = device.connectionLostHandler

	if isSecure {
		// tlsConfig := NewTlsConfig()
		// opts.SetTLSConfig(tlsConfig)
		// TODO:
	}

	device.client = mqtt.NewClient(opts)
	device.epZero = endpoints.NewZeroEndpoint(device.uid, "0", device.SendConfigs, device.sendMsg)
	return device

}

func (obj *MqttDevice) Connect() error {
	// obj.epZero = endpoints.NewZeroEndpoint(obj.uid, "0", obj.SendConfigs, obj.sendMsg)

	fmt.Print("Connecting\n")
	if obj.client == nil {
		fmt.Print("Client is nil")
		return errors.New("client is nil")
	}

	if token := obj.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	err := obj.subscribeTopic(fmt.Sprintf("%s/#", obj.uid))
	if err != nil {
		return err
	}

	err = obj.sendMsg(fmt.Sprintf("d/%s/0/announce", obj.uid), "")
	if err != nil {
		return err
	}

	return nil
}

func (obj *MqttDevice) Disconnect() {
	obj.client.Disconnect(250)
	fmt.Println("client disconnected")
}

func (obj *MqttDevice) RunForever(quit chan int) {
	for {
		<-quit
		fmt.Println("quit")
		break
	}
	fmt.Println("run forever done")
}

func (obj *MqttDevice) onMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	epID := parseEndpointID(msg.Topic())
	fmt.Printf("PARSED EP %s\n", epID)
	if epID == "" {
		return
	}
	if epID == "0" {
		obj.epZero.HandleMessage(string(msg.Payload()))
		return
	}

	if enp, found := obj.endpoints[epID]; found {
		enp.HandleMessage(parseCommand(msg.Topic()), string(msg.Payload()))
	} else {
		fmt.Println("Endpoint not found in map")
	}

}

func (obj *MqttDevice) onConnectHandler(client mqtt.Client) {
	fmt.Printf("Client with ID %s connected\n", obj.uid)
}

func (obj *MqttDevice) connectionLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Client with ID %s lost connection\n", obj.uid)
}

func (obj *MqttDevice) AddEndpoint(enp endpoints.Endpoint) {
	enp.SetOwnerID(obj.uid)
	enp.RegisterSendMsgCb(obj.sendMsg)
	obj.endpoints[enp.GetID()] = enp
}

func (obj *MqttDevice) SendConfigs() {

	fmt.Printf("Sending configs\n")
	// send zero endpoint config
	obj.epZero.SendConfig(len(obj.endpoints))
	// send endpoints configs
	for _, value := range obj.endpoints {
		value.SendConfig()
	}
}

func (obj *MqttDevice) subscribeTopic(topic string) error {
	token := obj.client.Subscribe(topic, 1, nil)
	success := token.WaitTimeout(time.Second * 2)
	if !success {
		return errors.New("failed to subscribe")
	}
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
	return nil
}

func (obj *MqttDevice) sendMsg(topic string, msg string) error {
	fmt.Printf("Sending MSG Topic: %s Message: %s\n", topic, msg)
	token := obj.client.Publish(topic, 1, false, msg)
	success := token.WaitTimeout(time.Second * 2)
	if !success {
		return errors.New("failed to announce")
	}
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}
