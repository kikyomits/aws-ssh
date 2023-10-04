# aws-ssh

## Overview

`aws-ssh` allows developers to execute commands inside Amazon Elastic Container Service (ECS) containers and set up port forwarding to remote hosts or ports. It simplifies the process of interacting with containers running in ECS clusters, making it easy to manage and troubleshoot containerized applications.

## Features

- Execute custom commands inside ECS containers.
- Forward network traffic from a local port to a remote host and port securely.
- Works seamlessly with ECS clusters and tasks.

## Prerequisites

TBD

## Installation

TBD

## Usage

### Execute a Command in an ECS Container

To execute a command inside an ECS container, use the following command:

```sh
your-cli -c <ECS_CLUSTER_NAME> -t <ECS_TASK_NAME> -cmd "<YOUR_COMMAND>"
```

- `-c, --cluster`: Specifies the ECS cluster name where the container is running.
- `-t, --task`: Specifies the ECS task or service name containing the container.
- `-cmd, --command`: Specifies the command to execute inside the container.

### Set Up Port Forwarding

To set up port forwarding from a local port to a remote host and port, use the following command:

```sh
your-cli -lp <LOCAL_PORT> -rh <REMOTE_HOST> -rp <REMOTE_PORT>
```

- `-lp, --local-port`: Specifies the local port to forward traffic to/from the container.
- `-rh, --remote-host`: Specifies the remote host to forward traffic to/from.
- `-rp, --remote-port`: Specifies the remote port to forward traffic to/from.

## Examples

Here are some examples of how to use the CLI tool:

```sh
# Execute a custom commands in an ECS container
your-cli -c my-cluster -t my-task -cmd "echo Hello, World!"

# Set up port forwarding from local port 8080 to a remote web server
your-cli -lp 8080 -rh remote-host -rp 80
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Special thanks to the open-source community for providing valuable tools and libraries.

---

Feel free to customize this README template to include specific details about your CLI tool, such as installation instructions, usage examples, and any additional features or options it offers.