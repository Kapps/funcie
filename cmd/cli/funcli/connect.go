package funcli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"strings"
)

type ConnectConfig struct {
	RemoteHost string `arg:"--remote-host,-r" help:"Override the remote host to connect to instead of the Redis default."`
	RemotePort int    `arg:"--remote-port,-p" help:"Override the remote port to bind to." default:"6379"`
	LocalPort  int    `arg:"--local-port,-l" help:"Override the local port to bind to." default:"6379"`
}

type ConnectCommand struct {
	cliConfig           *CliConfig
	configStore         ConfigStore
	connectClient       SsmClient
	tunneller           Tunneller
	connectivityService ConnectivityService
}

func NewConnectCommand(
	cliConfig *CliConfig,
	configStore ConfigStore,
	connectClient SsmClient,
	tunneller Tunneller,
	connectivityService ConnectivityService,
) *ConnectCommand {

	return &ConnectCommand{
		cliConfig:           cliConfig,
		configStore:         configStore,
		connectClient:       connectClient,
		tunneller:           tunneller,
		connectivityService: connectivityService,
	}
}

func (c *ConnectCommand) Run(ctx context.Context) error {
	conf := c.cliConfig.ConnectConfig

	if conf.RemoteHost == "" {
		host, err := c.configStore.GetConfigValue(ctx, "redis_host")
		if err != nil {
			return fmt.Errorf("failed to get Redis host: %w", err)
		}
		conf.RemoteHost = strings.Split(host, ":")[0]
	}

	ssmEndpoint := fmt.Sprintf("https://ssm.%v.amazonaws.com", c.cliConfig.Region)

	for ctx.Err() == nil {
		if err := c.connectivityService.WaitForConnectivity(ctx, ssmEndpoint); err != nil {
			return fmt.Errorf("failed to wait for connectivity: %w", err)
		}

		if err := c.startTunnel(ctx); err != nil {
			return fmt.Errorf("failed to start tunnel: %w", err)
		}
	}

	return nil
}

func (c *ConnectCommand) startTunnel(ctx context.Context) error {
	conf := c.cliConfig.ConnectConfig

	instanceId, err := c.configStore.GetConfigValue(ctx, "bastion_instance_id")
	if err != nil {
		return fmt.Errorf("failed to get instance ID: %w", err)
	}

	sess, err := c.connectClient.StartSession(ctx, &ssm.StartSessionInput{
		Target:       aws.String(instanceId),
		DocumentName: aws.String("AWS-StartPortForwardingSessionToRemoteHost"),
		Parameters: map[string][]string{
			"portNumber":      {fmt.Sprintf("%v", conf.RemotePort)},
			"localPortNumber": {fmt.Sprintf("%v", conf.LocalPort)},
			"host":            {conf.RemoteHost},
		},
		Reason: nil,
	})
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}

	opts := &SsmTunnellerOptions{
		Output:     sess,
		InstanceID: instanceId,
	}
	err = c.tunneller.OpenTunnel(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to open tunnel: %w", err)
	}

	_, err = c.connectClient.TerminateSession(ctx, &ssm.TerminateSessionInput{
		SessionId: sess.SessionId,
	})
	if err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}

	return nil
}
