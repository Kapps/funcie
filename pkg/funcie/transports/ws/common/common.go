package common

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

type ClientToServerMessage struct {
	RequestType string           `json:"requestType"`
	Application string           `json:"channel"`
	Response    *funcie.Response `json:"response"`
}

type ServerToClientMessage struct {
	RequestType string            `json:"requestType"`
	Message     *messages.Message `json:"message"`
}

const ClientToServerMessageRequestTypeSubscribe = "s"
const ClientToServerMessageRequestTypeUnsubscribe = "u"
const ClientToServerMessageRequestTypeResponse = "rs"

const ServerToClientMessageRequestTypeRequest = "rq"
