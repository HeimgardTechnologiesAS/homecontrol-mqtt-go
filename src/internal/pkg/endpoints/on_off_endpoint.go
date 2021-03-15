package endpoints

import (
	"errors"
	"fmt"
)

const epType = "pwr"
const reportingTime = "60"

type OnOffEndpoint struct {
	ownerID          string
	id               string
	state            string
	onStateChangedCb func(ep Endpoint, cmd string, state string)
	sendFeedbackCb   func(topic string, msg string) error
}

func NewOnOffEndpoint(
	epId string,
	onStateChange func(ep Endpoint, cmd string, state string),
) *OnOffEndpoint {
	return &OnOffEndpoint{
		id:               epId,
		state:            "0",
		onStateChangedCb: onStateChange,
	}
}

func (obj *OnOffEndpoint) SetOwnerID(id string) {
	obj.ownerID = id
}

func (obj *OnOffEndpoint) GetOwnerID() string {
	return obj.ownerID
}

func (obj *OnOffEndpoint) GetID() string {
	return obj.id
}

func (obj *OnOffEndpoint) HandleMessage(cmd string, msg string) {
	obj.state = msg
	if obj.onStateChangedCb != nil {
		obj.onStateChangedCb(obj, cmd, msg)
	}
}

func (obj *OnOffEndpoint) RegisterOnStateChangedCb(cb func(ep Endpoint, cmd string, state string)) {
	obj.onStateChangedCb = cb
}

func (obj *OnOffEndpoint) RegisterSendMsgCb(cb func(topic string, msg string) error) {
	obj.sendFeedbackCb = cb
}

func (obj *OnOffEndpoint) SendConfig() {
	if obj.sendFeedbackCb != nil {
		obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/conf", obj.ownerID, obj.id), fmt.Sprintf("e=%s;r=%s;", epType, reportingTime))
	}
}

func (obj *OnOffEndpoint) SendFeedbackMessage(cmd string, msg string) error {
	if msg != "1" && msg != "0" {
		return errors.New("unsupported message type for ON/OFF endpoint")
	}
	if obj.sendFeedbackCb != nil {
		obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/%s", obj.ownerID, obj.id, cmd), msg)
		return nil
	}
	return errors.New("callback not set")
}
