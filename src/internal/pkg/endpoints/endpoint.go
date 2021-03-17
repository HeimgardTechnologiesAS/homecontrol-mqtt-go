package endpoints

type Endpoint interface {
	SetOwnerID(id string)
	GetOwnerID() string
	GetID() string
	HandleMessage(cmd string, msg string)
	RegisterOnStateChangedCb(cb func(ep Endpoint, cmd string, state string))
	RegisterSendMsgCb(cb func(topic string, msg string) error)
	SendConfig()
	SendFeedbackMessage(cmd string, msg string) error
	SendStatus()
}
