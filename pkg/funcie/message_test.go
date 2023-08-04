package funcie_test

import (
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalPayload(t *testing.T) {
	regPayload := messages.NewRegistrationRequestPayload("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))
	serializedPayload := funcie.MustSerialize(regPayload)

	marshaledMessage := funcie.NewMessage("name", messages.MessageKindRegister, serializedPayload)

	unmarshaledMessage, err := funcie.UnmarshalMessagePayload[messages.RegistrationMessage](marshaledMessage)
	require.NoError(t, err)

	require.Equal(t, *regPayload, unmarshaledMessage.Payload)
}

func TestMarshalPayload(t *testing.T) {
	regPayload := messages.NewRegistrationRequestPayload("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))

	unmarshaledMessage := funcie.NewMessageWithPayload("name", messages.MessageKindRegister, regPayload)

	marshaledMessage, err := funcie.MarshalMessagePayload(*unmarshaledMessage)
	require.NoError(t, err)

	require.Equal(t, json.RawMessage(funcie.MustSerialize(regPayload)), marshaledMessage.Payload)
}

func TestMarshalMessagePayload(t *testing.T) {
	json := json.RawMessage(`{"version":"2.0","routeKey":"$default","rawPath":"/","rawQueryString":"name=Bob","headers":{"sec-fetch-mode":"navigate","x-amzn-tls-version":"TLSv1.2","sec-fetch-site":"none","accept-language":"en-CA,en-US;q=0.7,en;q=0.3","x-forwarded-proto":"https","x-forwarded-port":"443","x-forwarded-for":"74.12.30.203","sec-fetch-user":"?1","accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8","x-amzn-tls-cipher-suite":"ECDHE-RSA-AES128-GCM-SHA256","x-amzn-trace-id":"Root=1-64cc5cce-3e320ebf3cf60fc74f112635","host":"osowwuc56g7cpjjeqrcfjzfmxa0lwxuk.lambda-url.us-east-1.on.aws","upgrade-insecure-requests":"1","accept-encoding":"gzip, deflate, br","sec-fetch-dest":"document","user-agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/114.0"},"queryStringParameters":{"name":"Bob"},"requestContext":{"accountId":"anonymous","apiId":"osowwuc56g7cpjjeqrcfjzfmxa0lwxuk","domainName":"osowwuc56g7cpjjeqrcfjzfmxa0lwxuk.lambda-url.us-east-1.on.aws","domainPrefix":"osowwuc56g7cpjjeqrcfjzfmxa0lwxuk","http":{"method":"GET","path":"/","protocol":"HTTP/1.1","sourceIp":"74.12.30.203","userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/114.0"},"requestId":"ef43a940-5b1f-4cbb-b70a-f08cc9e87586","routeKey":"$default","stage":"$default","time":"04/Aug/2023:02:05:02 +0000","timeEpoch":1691114702522},"isBase64Encoded":false}`)
	payload := messages.NewForwardRequestPayload(json)

	message := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, *payload)

	marshal, err := funcie.MarshalMessagePayload(*message)
	require.NoError(t, err)

	unmarshaled, err := funcie.UnmarshalMessagePayload[messages.ForwardRequestMessage](marshal)
	require.NoError(t, err)

	require.Equal(t, message, unmarshaled)
}
