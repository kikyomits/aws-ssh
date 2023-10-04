package sessions

import (
	"github.com/aws/session-manager-plugin/src/log"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"  // this is required to register portSession
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession" // this is required to register shellSession
)

func NewAWSPluginSession() *PluginSession {
	return &PluginSession{}
}

type PluginSession struct {
}

func (p *PluginSession) Execute(in PluginSessionInput) error {
	ses := in.toSession()
	return ses.Execute(log.Logger(false, in.ClientId))
}
