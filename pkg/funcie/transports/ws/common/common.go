package common

import (
	"github.com/Kapps/funcie/pkg/funcie"
)

type ClientToServerMessage struct {
	RequestType string           `json:"requestType"`
	Application string           `json:"channel"`
	Response    *funcie.Response `json:"response"`
}

type ServerToClientMessage struct {
	RequestType string          `json:"requestType"`
	Message     *funcie.Message `json:"message"`
}

const ClientToServerMessageRequestTypeSubscribe = "s"
const ClientToServerMessageRequestTypeUnsubscribe = "u"
const ClientToServerMessageRequestTypeResponse = "rs"

const ServerToClientMessageRequestTypeRequest = "rq"
