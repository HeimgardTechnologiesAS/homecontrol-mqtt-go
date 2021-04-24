package endpoints

import (
	"errors"
	"fmt"
	c "homecontrol-mqtt-go/internal/pkg/commands"
	"log"
)

// endpoint struct should be embedded to all child endpoints.
// It defines all methods that one endpint must have
type endpoint struct {
	epType           string
	reportingTime    string
	ownerID          string
	id               string
	name             string
	sendFeedbackCb   func(topic string, msg string) error
	onStateChangedCb func(ep Endpoint, cmd string, state string)
	commands         map[string]c.Command
}

// newEndpoint constructs a new endpoint
func newEndpoint(
	epType string,
	epReportingTime string,
	epId string,
	epName string,
	epOnStateChange func(ep Endpoint, cmd string, state string),
	epCommands map[string]c.Command,
) *endpoint {
	return &endpoint{
		epType:           epType,
		reportingTime:    epReportingTime,
		id:               epId,
		name:             epName,
		onStateChangedCb: epOnStateChange,
		commands:         epCommands,
	}
}

// SetOwnerID sets endpoint's owner ID (device ID)
func (obj *endpoint) SetOwnerID(id string) {
	obj.ownerID = id
}

// GetOwnerID returns owner ID
func (obj *endpoint) GetOwnerID() string {
	return obj.ownerID
}

// GetID returns endpoint ID
func (obj *endpoint) GetID() string {
	return obj.id
}

// RegisterSendMsgCb registers callback function that is called when message from endpoint must be sent to HC GW
func (obj *endpoint) RegisterSendMsgCb(cb func(topic string, msg string) error) {
	obj.sendFeedbackCb = cb
}

// RegisterOnStateChangedCb registers callback funkcion that is called when endpoint state is changed
func (obj *endpoint) RegisterOnStateChangedCb(cb func(ep Endpoint, cmd string, state string)) {
	obj.onStateChangedCb = cb
}

// HandleMessage handles incoming message
func (obj *endpoint) HandleMessage(cmd string, msg string) {
	val, ok := obj.commands[cmd]
	if !ok {
		log.Printf("received command not supported %s, ep ID: %s", cmd, obj.id)
		return
	}
	val.SetState(msg)
	if obj.onStateChangedCb != nil {
		obj.onStateChangedCb(obj, cmd, msg)
	}
}

// SendFeedbackMessage sends feedback message to HC GW when some of the commands change state
func (obj *endpoint) SendFeedbackMessage(cmd string, msg string) error {
	val, ok := obj.commands[cmd]
	if !ok {
		return fmt.Errorf("unsupported command type provided: %s", cmd)
	}
	val.SetState(msg)
	if obj.sendFeedbackCb != nil {
		obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/%s", obj.ownerID, obj.id, cmd), msg)
		return nil
	}
	return errors.New("callback not set")
}

// SendConfig sends endpoint config to HC GW
func (obj *endpoint) SendConfig() {
	if obj.sendFeedbackCb != nil {
		cfg := fmt.Sprintf("e=%s;r=%s;", obj.epType, obj.reportingTime)
		if obj.name != "" {
			cfg = fmt.Sprintf("%sname=%s;", cfg, obj.name)
		}
		obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/conf", obj.ownerID, obj.id), cfg)
	}
}

// SendStatus sends current endpoint commands status
func (obj *endpoint) SendStatus() {
	log.Printf("Not implemented. Endpoint ID: %s, Endpoint type: %s", obj.id, obj.epType)
}
