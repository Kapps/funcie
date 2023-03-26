package common

type ClientToServerMessage struct {
	Channel     string `json:"channel"`
	RequestType string `json:"requestType"`
}

type ServerToClientMessage struct {
	Channel     string `json:"channel"`
	Payload     string `json:"payload"`
	RequestType string `json:"requestType"`
}

const ClientToServerMessageRequestTypeSubscribe = "s"
const ClientToServerMessageRequestTypeUnsubscribe = "u"

//const ServerToClientMessageRequestTypeInvoke = "invoke"
