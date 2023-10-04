package sessions

import (
	"aws-ssh/internal/services"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/google/uuid"
)

type SSMSessionManager struct {
	plugin Plugin
	ssm    SSMClient
	region string
}

func NewSSMSessionManager(plugin Plugin, ssm SSMClient, region string) (*SSMSessionManager, error) {
	return &SSMSessionManager{
		plugin: plugin,
		ssm:    ssm,
		region: region,
	}, nil
}

// PortForwardingToRemoteHostSession starts a port forwarding sessions using the PortForwardingInput parameters to
// configure the sessions.  The aws.Config parameter will be used to call the AWS SSM StartSession
// API, which is used as part of establishing the websocket communication channel.
func (c *SSMSessionManager) PortForwardingToRemoteHostSession(opts *PortForwardingToRemoteInput) error {
	in := &ssm.StartSessionInput{
		DocumentName: aws.String(services.SSMDocumentAWSStartPortForwardingSessionToRemoteHost),
		Target:       aws.String(opts.Target),
		Parameters: map[string][]string{
			"localPortNumber": {opts.LocalPort},
			"portNumber":      {opts.RemotePort},
			"host":            {opts.Host},
		},
	}
	return c.execute(in)
}

func (c *SSMSessionManager) PortForwardingSession(opts *PortForwardingInput) error {
	in := &ssm.StartSessionInput{
		DocumentName: aws.String(services.SSMDocumentAWSStartPortForwardingSession),
		Target:       aws.String(opts.Target),
		Parameters: map[string][]string{
			"localPortNumber": {opts.LocalPort},
			"portNumber":      {opts.RemotePort},
		},
	}
	return c.execute(in)
}

func (c *SSMSessionManager) execute(input *ssm.StartSessionInput) error {
	out, err := c.ssm.StartSession(context.Background(), input)
	if err != nil {
		return err
	}
	ep, err := ssm.NewDefaultEndpointResolver().ResolveEndpoint(c.region, ssm.EndpointResolverOptions{})
	if err != nil {
		return err
	}
	pluginInput := PluginSessionInput{
		ClientId:    uuid.NewString(),
		DataChannel: &datachannel.DataChannel{},
		Endpoint:    ep.URL,
		SessionId:   *out.SessionId,
		StreamUrl:   *out.StreamUrl,
		TargetId:    *input.Target,
		TokenValue:  *out.TokenValue,
	}
	return c.plugin.Execute(pluginInput)
}
