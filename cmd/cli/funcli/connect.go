package funcli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type ConnectCommand struct {
	RemoteHost string `arg:"--remote-host,-r" help:"Override the remote host to connect to instead of the Redis default."`
	RemotePort int    `arg:"--remote-port,-p" help:"Override the remote port to bind to." default:"6379"`
	LocalPort  int    `arg:"--local-port,-l" help:"Override the local port to bind to." default:"6379"`

	configStore   ConfigStore
	connectClient SsmConnectClient
	tunneller     Tunneller
}

type SsmConnectClient interface {
	StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)
	TerminateSession(ctx context.Context, params *ssm.TerminateSessionInput, optFns ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error)
}

func NewConnectCommand(configStore ConfigStore, connectClient SsmConnectClient, tunneller Tunneller) *ConnectCommand {
	return &ConnectCommand{
		configStore:   configStore,
		connectClient: connectClient,
		tunneller:     tunneller,
	}
}

func (c *ConnectCommand) Run(ctx context.Context) error {
	if c.RemoteHost == "" {
		host, err := c.configStore.GetConfigValue(ctx, "redis_host")
		if err != nil {
			return fmt.Errorf("failed to get Redis host: %w", err)
		}
		c.RemoteHost = host
	}

	if err := c.startTunnel(ctx); err != nil {
		return fmt.Errorf("failed to start tunnel: %w", err)
	}

	return nil
}

func (c *ConnectCommand) startTunnel(ctx context.Context) error {
	instanceId, err := c.configStore.GetConfigValue(ctx, "bastion_instance_id")
	if err != nil {
		return fmt.Errorf("failed to get instance ID: %w", err)
	}

	sess, err := c.connectClient.StartSession(ctx, &ssm.StartSessionInput{
		Target:       aws.String(instanceId),
		DocumentName: aws.String("AWS-StartPortForwardingSessionToRemoteHost"),
		Parameters: map[string][]string{
			"portNumber":      {fmt.Sprintf("%v", c.RemotePort)},
			"localPortNumber": {fmt.Sprintf("%v", c.LocalPort)},
			"host":            {c.RemoteHost},
		},
		Reason: nil,
	})
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}

	fmt.Printf("%+v", sess.StreamUrl)

	fmt.Println("Session started:", *sess.SessionId, "Stream URL:", *sess.StreamUrl)
	opts := &TunnelOptions{
		Headers: map[string][]string{
			"X-Amz-Security-Token": {*sess.TokenValue},
		},
	}
	err = c.tunneller.OpenTunnel(ctx, *sess.StreamUrl, c.LocalPort, opts)
	if err != nil {
		return fmt.Errorf("failed to open tunnel: %w", err)
	}

	/*_, err = c.connectClient.TerminateSession(ctx, &ssm.TerminateSessionInput{
		SessionId: sess.SessionId,
	})
	if err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}*/

	return nil
}
