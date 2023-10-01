package client

type Client interface {
	httpForward()
	sshTunnel()
	sshTunnelRemote()
}
