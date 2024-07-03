package internal

import (
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"net/url"
	"os"
)

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

func RequireUrlEnv(name string, purpose ConfigPurpose) url.URL {
	value := RequiredEnv(name, purpose)
	parsedUrl, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s %s: %s", name, value, err))
	}
	return *parsedUrl
}

func RequiredEnv(name string, purpose ConfigPurpose) string {
	value := os.Getenv(name)
	if value == "" {
		currPurpose := GetConfigPurpose()
		purposeMatches := purpose == ConfigPurposeAny || currPurpose == purpose
		if purposeMatches {
			panic(fmt.Sprintf("required environment variable %s not set", name))
		}
	}
	return value
}

func OptionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func OptionalUrlEnv(name string, defaultValue string) url.URL {
	value := OptionalEnv(name, defaultValue)
	parsedUrl, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s %s: %s", name, value, err))
	}
	return *parsedUrl
}
