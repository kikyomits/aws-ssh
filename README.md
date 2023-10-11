# aws-ssh

## Overview

`aws-ssh` helps developers to execute commands inside Amazon Elastic Container Service (ECS) containers and set up port forwarding to remote hosts or ports. It simplifies the process of interacting with containers running in ECS clusters, making it easy to manage and troubleshoot containerized applications.

```
Usage:
  aws-ssh [command]

Available Commands:
  completion       Generate the autocompletion script for the specified shell
  ecs-port-forward Port forwarding for AWS ECS tasks.
  help             Help about any command

Flags:
  -h, --help             help for aws-ssh
      --profile string   AWS Profile
      --region string    AWS Region
  -v, --verbose          enable verbose logging

Use "aws-ssh [command] --help" for more information about a command.
```

## Features

- Execute custom commands inside ECS containers.
- Forward network traffic from a local port to a remote host and port securely.
- Works seamlessly with ECS clusters and tasks.

## Installation

### homebrew

You can install `aws-ssh` via [homebrew](https://brew.sh/).

```shell
brew tap kikyomits/aws-ssh https://github.com/kikyomits/aws-ssh
brew install aws-ssh
```

### Release Page

If you don't have `brew` in your machine, you can download the cli from [Releases](https://github.com/kikyomits/aws-ssh/releases).
Please find the executable file for your platform.

## Usage

### Set Up Port Forwarding

To set up port forwarding from a local port to a remote host and port, use the following command:


```
Usage: aws-ssh ecs_ssh-port-forward [ECSPortForwardOptions]
        Forward localPort port to localPort port on Task
        Forward localPort port to a remote host/port accessible from Task

Usage:
  aws-ssh ecs-port-forward [flags]

Examples:
aws-ssh ecs-port-forward --Cluster CLUSTER_NAME --Local LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT --Task TASK_ID

Flags:
  -c, --cluster string     ECS Cluster Name
  -n, --container string   Container name. Required if Task is running more than one Container
  -h, --help               help for ecs-port-forward
  -L, --local string       LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT Forward a Local port to a remote address/port
  -s, --service string     ECS Service Name. If provided, it will search the ECS Service and try to access to an Active Task. Either of Service or Task must be provided.
  -t, --task string        ECS Task ID. Either of Service or Task must be provided.

Global Flags:
      --profile string   AWS Profile
      --region string    AWS Region
  -v, --verbose          enable verbose logging
```

### Execute a Command in an ECS Container

To execute a command inside an ECS container, use the following command:

```sh
TBD
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
