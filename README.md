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

```

## Features

- Execute custom commands inside ECS containers.
- Forward network traffic from a local port to a remote host and port securely.
- Works seamlessly with ECS clusters and tasks.

## Installation

TBD

## Usage

### Execute a Command in an ECS Container

To execute a command inside an ECS container, use the following command:

```sh
TBD
```


### Set Up Port Forwarding

To set up port forwarding from a local port to a remote host and port, use the following command:


```sh
Usage: aws-ssh ecs_ssh-port-forward [ECSPortForwardOptions]
        Forward localPort port to localPort port on task
        Forward localPort port to a remote host/port accessible from task

Usage:
  aws-ssh ecs_ssh-port-forward [flags]

Examples:
aws-ssh ecs_ssh-port-forward --cluster CLUSTER_NAME --local LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT --task TASK_ID

Flags:
  -c, --cluster string     ECS Cluster Name
  -n, --container string   Container name. Required if task is running more than one container
  -h, --help               help for ecs_ssh-port-forward
  -L, --local string       LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT Forward a local port to a remote address/port
  -s, --ecs_ssh string    ECS Service Name. If provided, it will search the ECS Service and try to access to an Active task. Either of Service or Task must be provided.
  -t, --task string        ECS Task ID. Either of Service or Task must be provided.

Global Flags:
      --profile string   AWS Profile
      --region string    AWS Region
  -v, --verbose          enable verbose logging
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
