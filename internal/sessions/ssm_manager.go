package sessions

import (
	"aws-ssh/internal/ecs_ssh"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type SSMSessionManager struct {
	plugin Plugin
	ssm    SSM
}

func NewSSMSessionManager(plugin Plugin, ssm SSM) *SSMSessionManager {
	return &SSMSessionManager{
		plugin: plugin,
		ssm:    ssm,
	}
}

// PortForwardingToRemoteHostSession starts a port forwarding sessions using the PortForwardingToRemoteInput parameters to
// configure the sessions.
func (c *SSMSessionManager) PortForwardingToRemoteHostSession(opts *PortForwardingToRemoteInput) error {
	log.Infof("setting up tunnel: 127.0.0.1:%s -> %s:%s", opts.LocalPort, opts.RemoteHost, opts.RemotePort)
	in := &ssm.StartSessionInput{
		DocumentName: aws.String(ecs_ssh.SSMDocumentAWSStartPortForwardingSessionToRemoteHost),
		Target:       aws.String(opts.Target),
		Parameters: map[string][]string{
			"localPortNumber": {opts.LocalPort},
			"portNumber":      {opts.RemotePort},
			"host":            {opts.RemoteHost},
		},
	}
	return c.execute(in, opts.Region)
}

// PortForwardingSession starts a port forwarding sessions using the PortForwardingInput parameters to
// configure the sessions.
func (c *SSMSessionManager) PortForwardingSession(opts *PortForwardingInput) error {
	log.Infof("setting up tunnel: 127.0.0.1:%s -> :%s", opts.LocalPort, opts.RemotePort)
	in := &ssm.StartSessionInput{
		DocumentName: aws.String(ecs_ssh.SSMDocumentAWSStartPortForwardingSession),
		Target:       aws.String(opts.Target),
		Parameters: map[string][]string{
			"localPortNumber": {opts.LocalPort},
			"portNumber":      {opts.RemotePort},
		},
	}
	return c.execute(in, opts.Region)
}

func (c *SSMSessionManager) ExecSession(opts *ExecInput) error {
	in := &ssm.StartSessionInput{
		DocumentName: aws.String(ecs_ssh.SSMStartInteractiveCommand),
		Target:       aws.String(opts.Target),
		Parameters: map[string][]string{
			"command": {opts.Command},
		},
	}
	return c.execute(in, opts.Region)
}

func (c *SSMSessionManager) execute(input *ssm.StartSessionInput, region string) error {
	out, err := c.ssm.StartSession(context.Background(), input)
	if err != nil {
		return err
	}

	log.WithFields(map[string]interface{}{
		"region":    region,
		"sessionId": aws.ToString(out.SessionId),
		"streamUrl": aws.ToString(out.StreamUrl),
	}).Debugf("connecting to stream via AWS SessionManager plugin")

	ep, err := ssm.NewDefaultEndpointResolverV2().ResolveEndpoint(context.Background(), ssm.EndpointParameters{Region: aws.String(region)})
	if err != nil {
		log.WithError(err).Errorf("failed to resolve endpoint for plugin")
		return err
	}

	pluginInput := PluginSessionInput{
		ClientId:    uuid.NewString(),
		DataChannel: &datachannel.DataChannel{},
		Endpoint:    ep.URI.String(),
		SessionId:   *out.SessionId,
		StreamUrl:   *out.StreamUrl,
		TargetId:    *input.Target,
		TokenValue:  *out.TokenValue,
	}

	err = c.plugin.Execute(pluginInput)
	return err
}
