package devices

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/HomeControlAS/homecontrol-mqtt-go/endpoints"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TLS_PORT  = 8883 // secure TLS port
	TCP_PORT  = 1883 // not secure port
	zeroEnpID = "0"  // zero endpoint ID, must be 0
)

// MqttDevice struct that defines all method needed to connect and communicate with HC GW
type MqttDevice struct {
	uid       string
	name      string
	client    mqtt.Client
	epZero    *endpoints.ZeroEndpoint
	endpoints map[string]endpoints.Endpoint
	quitC     chan error
}

// NewMqttDevice constructs Mqtt device
func NewMqttDevice(
	hostname string,
	uid string,
	username string,
	password string,
	isSecure bool,
	deviceName string,
) (Device, error) {

	device := &MqttDevice{
		uid:       uid,
		name:      deviceName,
		endpoints: make(map[string]endpoints.Endpoint),
	}

	opts := mqtt.NewClientOptions()

	if isSecure {
		opts.AddBroker(fmt.Sprintf("ssl://%s:%d", hostname, TLS_PORT))
	} else {
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", hostname, TCP_PORT))
	}
	opts.SetClientID(uid)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetWill(fmt.Sprintf("d/%s/0/online", uid), fmt.Sprintf("d/%s/0/offline", uid), 1, false)

	opts.SetDefaultPublishHandler(device.onMessageHandler)
	opts.OnConnect = device.onConnectHandler
	opts.OnConnectionLost = device.connectionLostHandler

	if isSecure {
		tlsConfig, err := fetchCertificate(hostname)
		if err != nil {
			return nil, err
		}
		opts.SetTLSConfig(tlsConfig)
	}

	device.client = mqtt.NewClient(opts)
	device.epZero = endpoints.NewZeroEndpoint(device.uid, zeroEnpID, device.SendConfigs, device.sendMsg)

	return device, nil
}

// fetchCertificate fetches tls cert from HC gateway
func fetchCertificate(hostname string) (*tls.Config, error) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, TLS_PORT))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to host: %s", err.Error())
	}

	client := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		return nil, fmt.Errorf("fetching certificate client handshake failed: %s", err.Error())
	}

	cert := client.ConnectionState().PeerCertificates[0]
	certpool := x509.NewCertPool()
	certpool.AddCert(cert)

	return &tls.Config{
		RootCAs:            certpool,
		MaxVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}, nil
}

// Connect connects Device object to HC gateway
func (obj *MqttDevice) Connect() error {
	if obj.client == nil {
		return errors.New("client is not initialized")
	}

	token := obj.client.Connect()
	success := token.WaitTimeout(time.Second * 5)
	if !success {
		return fmt.Errorf("timeout ocurred while trying to connect to client")
	}
	if token.Error() != nil {
		return fmt.Errorf("got error while connecting to client %s", token.Error())
	}

	err := obj.subscribeTopic(fmt.Sprintf("%s/#", obj.uid))
	if err != nil {
		return err
	}

	err = obj.subscribeTopic("broadcast")
	if err != nil {
		return err
	}

	err = obj.announce()
	if err != nil {
		return err
	}

	return nil
}

// Disconnect disconnects Device from HC gateway
func (obj *MqttDevice) Disconnect() {
	obj.sendMsg(fmt.Sprintf("d/%s/0/offline", obj.uid), "")
	obj.client.Disconnect(250)
}

// RunForever runs infinite loop if MQTT Device should listen forever
func (obj *MqttDevice) RunForever(quitC chan error) error {
	obj.quitC = quitC
	for {
		err := <-quitC
		return err
	}
}

// AddEndpoint adds new endpoint to MQTT Device
func (obj *MqttDevice) AddEndpoint(enp endpoints.Endpoint) {
	enp.SetOwnerID(obj.uid)
	enp.RegisterSendMsgCb(obj.sendMsg)
	obj.endpoints[enp.GetID()] = enp
}

// SendConfigs sends config for each endpoint
func (obj *MqttDevice) SendConfigs() {
	// send zero endpoint config
	obj.epZero.SendConfig(len(obj.endpoints))
	// send endpoints configs
	for _, value := range obj.endpoints {
		value.SendConfig()
	}
}

// subscribeTopic subscribes MQTT to provided topic
func (obj *MqttDevice) subscribeTopic(topic string) error {
	token := obj.client.Subscribe(topic, 1, nil)
	success := token.WaitTimeout(time.Second * 2)
	if !success {
		return fmt.Errorf("failed to subscribe topic %s. Timeout occurred", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf("got error in token when subscribing %s", token.Error())
	}
	return nil
}

// announce announces MQTT device to HC gateway
func (obj *MqttDevice) announce() error {
	err := obj.sendMsg(fmt.Sprintf("d/%s/0/announce", obj.uid), obj.name)
	if err != nil {
		return fmt.Errorf("failed to announce device %s", err)
	}

	err = obj.sendMsg(fmt.Sprintf("d/%s/0/online", obj.uid), "")
	if err != nil {
		return fmt.Errorf("failed to send online device status %s", err)
	}

	for _, value := range obj.endpoints {
		err = value.SendStatus()
		if err != nil {
			return fmt.Errorf("failed to send status for endpoint %s, error: %s", value.GetID(), err)
		}
	}
	return nil
}

// sendMsg sends message to HC GW
func (obj *MqttDevice) sendMsg(topic string, msg string) error {
	token := obj.client.Publish(topic, 1, false, msg)
	success := token.WaitTimeout(time.Second * 2)
	if !success {
		return fmt.Errorf("failed to subscribe topic %s. Timeout occurred", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf(
			"got error in token when sending message: %s on topic: %s. %s",
			topic, msg, token.Error(),
		)
	}
	return nil
}

// onMessageHandler handles incoming mqtt messages
func (obj *MqttDevice) onMessageHandler(client mqtt.Client, msg mqtt.Message) {

	if strings.Contains(msg.Topic(), "broadcast") {
		if strings.Contains(string(msg.Payload()), "serverannounce") {
			obj.announce()
		}
	}

	epID := parseEndpointID(msg.Topic())

	if epID == zeroEnpID {
		obj.sendMsg(fmt.Sprintf("d/%s/0/online", obj.uid), "")
		obj.epZero.HandleMessage(string(msg.Payload()))
		return
	}

	if enp, found := obj.endpoints[epID]; found {
		enp.HandleMessage(parseCommand(msg.Topic()), string(msg.Payload()))
	}
}

// onConnectHandler handles MQTT connect event
func (obj *MqttDevice) onConnectHandler(client mqtt.Client) {
	// do nothing
}

// onConnectionLostHandler handles MQTT connection lost event
func (obj *MqttDevice) connectionLostHandler(client mqtt.Client, err error) {
	if obj.quitC != nil {
		obj.quitC <- fmt.Errorf("client with ID %s lost connection. Error: %s", obj.uid, err.Error())
	}
}

// parseEndpointID helper method to parse endpoint id from incoming message
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

// parseCommand helper method to parse command from incoming message
func parseCommand(token string) string {
	lDelimeterIn := strings.LastIndex(token, "/")
	if lDelimeterIn > -1 {
		return token[lDelimeterIn+1:]
	}
	return ""
}
