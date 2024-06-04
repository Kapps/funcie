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
)

// SsmTunnellerOptions contains the options for OpenTunnel using SSM.
type SsmTunnellerOptions struct {
	// Output is the output of the StartSession API call.
	Output *ssm.StartSessionOutput
	// InstanceID is the EC2 instance ID of the target instance.
	InstanceID string
}

// Tunneller is an interface for creating a tunnel to a remote host.
type Tunneller interface {
	// OpenTunnel starts a tunnel to a remote host using the provider-specific options.
	OpenTunnel(ctx context.Context, opts interface{}) error
}

type ssmTunnel struct {
	region string
}

// NewSsmTunneller creates a new WebhookTunnel.
func NewSsmTunneller(conf *CliConfig) Tunneller {
	return &ssmTunnel{
		region: conf.Region,
	}
}

func (t *ssmTunnel) OpenTunnel(ctx context.Context, opts interface{}) error {
	ssmOpts := opts.(*SsmTunnellerOptions)

	ep, err := ssm.NewDefaultEndpointResolver().ResolveEndpoint(t.region, ssm.EndpointResolverOptions{})
	if err != nil {
		return fmt.Errorf("failed to resolve endpoint: %w", err)
	}

	ssmSession := &session.Session{
		DataChannel: &datachannel.DataChannel{},
		SessionId:   *ssmOpts.Output.SessionId,
		StreamUrl:   *ssmOpts.Output.StreamUrl,
		TokenValue:  *ssmOpts.Output.TokenValue,
		Endpoint:    ep.URL,
		ClientId:    fmt.Sprintf("funcie-%v", uuid.NewString()),
		TargetId:    ssmOpts.InstanceID,
	}

	return ssmSession.Execute(log.Logger(false, ssmSession.ClientId))
}
