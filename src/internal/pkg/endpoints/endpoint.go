package endpoints

import (
	"errors"
	"fmt"
	c "homecontrol-mqtt-go/internal/pkg/commands"
	"log"
)

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

func (obj *endpoint) SetOwnerID(id string) {
	obj.ownerID = id
}

func (obj *endpoint) GetOwnerID() string {
	return obj.ownerID
}

func (obj *endpoint) GetID() string {
	return obj.id
}

func (obj *endpoint) RegisterSendMsgCb(cb func(topic string, msg string) error) {
	obj.sendFeedbackCb = cb
}

func (obj *endpoint) RegisterOnStateChangedCb(cb func(ep Endpoint, cmd string, state string)) {
	obj.onStateChangedCb = cb
}

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

func (obj *endpoint) SendConfig() {
	if obj.sendFeedbackCb != nil {
		cfg := fmt.Sprintf("e=%s;r=%s;", obj.epType, obj.reportingTime)
		if obj.name != "" {
			cfg = fmt.Sprintf("%sname=%s;", cfg, obj.name)
		}
		log.Printf("d/%s/%s/conf  %s", obj.ownerID, obj.id, cfg)
		obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/conf", obj.ownerID, obj.id), cfg)
	}
}

func (obj *endpoint) SendStatus() {
	log.Printf("Not implemented. Endpoint ID: %s, Endpoint type: %s", obj.id, obj.epType)
}
