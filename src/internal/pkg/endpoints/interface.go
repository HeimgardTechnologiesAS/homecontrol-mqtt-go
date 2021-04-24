package endpoints

// Endpoint interface collects all methods needed to handle device endpoint
type Endpoint interface {
	// SetOwnerID sets endpoint's owner ID (device ID)
	SetOwnerID(id string)
	// GetOwnerID returns owner ID
	GetOwnerID() string
	// GetID returns endpoint ID
	GetID() string
	// HandleMessage handles incoming message
	HandleMessage(cmd string, msg string)
	// RegisterOnStateChangedCb registers callback funkcion that is called when endpoint state is changed
	RegisterOnStateChangedCb(cb func(ep Endpoint, cmd string, state string))
	// RegisterSendMsgCb registers callback function that is called when message from endpoint must be sent to HC GW
	RegisterSendMsgCb(cb func(topic string, msg string) error)
	// SendConfig sends endpoint config to HC GW
	SendConfig()
	// SendFeedbackMessage sends feedback message to HC GW when some of the commands change state
	SendFeedbackMessage(cmd string, msg string) error
	// SendStatus sends current endpoint commands status
	SendStatus()
}
