package sessions

// PortForwardingInput configures the port forwarding sessions parameters.
// Target is the EC2 instance ID to establish the sessions with.
// RemotePort is the port on the EC2 instance to connect to.
// LocalPort is the port on the local host to listen to.  If not provided, a random port will be used.
type PortForwardingInput struct {
	Target     string
	RemotePort string
	LocalPort  string
}

// PortForwardingToRemoteInput configures the port forwarding sessions parameters.
// Target is the EC2 instance ID to establish the sessions with.
// RemotePort is the port on the EC2 instance to connect to.
// LocalPort is the port on the local host to listen to.  If not provided, a random port will be used.
// Host is the remote hostname
type PortForwardingToRemoteInput struct {
	Target     string
	RemotePort string
	LocalPort  string
	Host       string
}

type Manager interface {
	PortForwardingSession(in *PortForwardingInput) error
	PortForwardingToRemoteHostSession(in *PortForwardingToRemoteInput) error
}
