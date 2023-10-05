package internal

import "github.com/Kapps/funcie/pkg/funcie"

type ConfigPurpose int

var ConfigPurposeAny ConfigPurpose = 0
var ConfigPurposeClient ConfigPurpose = 1
var ConfigPurposeServer ConfigPurpose = 2

func GetConfigPurpose() ConfigPurpose {
	if funcie.IsRunningWithLambda() {
		return ConfigPurposeServer
	} else {
		return ConfigPurposeClient
	}
}
