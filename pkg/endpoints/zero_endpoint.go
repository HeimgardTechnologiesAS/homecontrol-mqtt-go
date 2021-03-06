package endpoints

import (
	"fmt"
	"strconv"
)

type ZeroEndpoint struct {
	ownerID        string
	id             string
	sendFeedbackCb func(topic string, msg string) error
	sendConfigsCb  func()
}

func NewZeroEndpoint(
	ownerID string,
	epId string,
	sendConfigsCb func(),
	sendFeedbackCb func(topic string, msg string) error,
) *ZeroEndpoint {
	return &ZeroEndpoint{
		ownerID:        ownerID,
		id:             epId,
		sendConfigsCb:  sendConfigsCb,
		sendFeedbackCb: sendFeedbackCb,
	}
}

func (obj *ZeroEndpoint) HandleMessage(msg string) {
	if obj.sendConfigsCb != nil {
		obj.sendConfigsCb()
	}
}

func (obj *ZeroEndpoint) SendConfig(enpCount int) error {
	if obj.sendFeedbackCb == nil {
		return fmt.Errorf("ZeroEndpoint sendFeedbackCb not set")
	}
	err := obj.sendFeedbackCb(fmt.Sprintf("d/%s/%s/conf", obj.ownerID, obj.id), strconv.Itoa(enpCount))
	if err != nil {
		return fmt.Errorf("ZeroEndpoint, failed to send config: %s", err)
	}
	return nil
}
