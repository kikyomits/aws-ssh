//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=manager.go -package=mock -destination=./mock/manager_mock.go
package sessions

// PortForwardingInput configures the port forwarding sessions parameters.
// Target is the EC2 instance ID to establish the sessions with.
// RemotePort is the port on the EC2 instance to connect to.
// LocalPort is the port on the local host to listen to.  If not provided, a random port will be used.
type PortForwardingInput struct {
	Region     string
	Target     string
	RemotePort string
	LocalPort  string
}

func NewPortForwardingInput(region, target, localPort, remotePort string) *PortForwardingInput {
	return &PortForwardingInput{
		Region:     region,
		Target:     target,
		LocalPort:  localPort,
		RemotePort: remotePort,
	}
}

// PortForwardingToRemoteInput configures the port forwarding sessions parameters.
// Target is the EC2 instance ID to establish the sessions with.
// RemotePort is the port on the EC2 instance to connect to.
// LocalPort is the port on the local host to listen to.  If not provided, a random port will be used.
// RemoteHost is the remote hostname
type PortForwardingToRemoteInput struct {
	Region     string
	Target     string
	LocalPort  string
	RemoteHost string
	RemotePort string
}

func NewPortForwardingToRemoteInput(region, target, localPort, remoteHost, remotePort string) *PortForwardingToRemoteInput {
	return &PortForwardingToRemoteInput{
		Region:     region,
		Target:     target,
		LocalPort:  localPort,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
	}
}

type ExecInput struct {
	Region  string
	Target  string
	Command string
}

func NewExecInput(region, target, command string) *ExecInput {
	return &ExecInput{
		Region:  region,
		Target:  target,
		Command: command,
	}
}

type Manager interface {
	ExecSession(in *ExecInput) error
	PortForwardingSession(in *PortForwardingInput) error
	PortForwardingToRemoteHostSession(in *PortForwardingToRemoteInput) error
}
