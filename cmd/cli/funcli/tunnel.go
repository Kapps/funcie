package funcli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/google/uuid"
	"net/http"
)

// TunnelOptions provides optional arguments for creating a network tunnel.
type TunnelOptions struct {
	// Headers is a map of headers to send with the connection request.
	Headers    http.Header
	Output     ssm.StartSessionOutput
	InstanceId string
}

// Tunneller is an interface for creating a tunnel to a remote host.
type Tunneller interface {
	// OpenTunnel starts a tunnel to a remote host on the given port locally.
	OpenTunnel(ctx context.Context, endpoint string, localPort int, opts *TunnelOptions) error
}

type webhookTunnel struct {
}

// NewWebhookTunneller creates a new WebhookTunnel.
func NewWebhookTunneller() Tunneller {
	return &webhookTunnel{}
}

func (t *webhookTunnel) OpenTunnel(ctx context.Context, endpoint string, localPort int, opts *TunnelOptions) error {
	ep, err := ssm.NewDefaultEndpointResolver().ResolveEndpoint("ca-central-1", ssm.EndpointResolverOptions{})
	if err != nil {
		return fmt.Errorf("failed to resolve endpoint: %w", err)
	}

	ssmSession := &session.Session{
		DataChannel: &datachannel.DataChannel{},
		SessionId:   *opts.Output.SessionId,
		StreamUrl:   *opts.Output.StreamUrl,

		TokenValue: *opts.Output.TokenValue,
		Endpoint:   ep.URL,
		ClientId:   fmt.Sprintf("funcie-%v", uuid.NewString()),
		TargetId:   opts.InstanceId,
	}

	return ssmSession.Execute(log.Logger(false, ssmSession.ClientId))
}
