//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=port_forward.go -package=mock -destination=./mock/port_forward_mock.go
package ecs_port_forward

import (
	"aws-ssh/internal/commands/factory"
	"aws-ssh/internal/commands/flags"
	"aws-ssh/internal/commands/prerun"
	"aws-ssh/internal/sessions"
	"aws-ssh/internal/validation"
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	clusterFlag   = "cluster"
	serviceFlag   = "services"
	taskFlag      = "task"
	localFlag     = "local"
	containerFlag = "container"
)

type ECSPortForwardOptions struct {
	cluster   string
	service   string
	task      string
	container string
	local     string
}

func (c *ECSPortForwardOptions) Validate(f factory.Factory, cmd *cobra.Command, args []string) error {
	if c.service == "" && c.task == "" {
		return validation.NewInvalidInputError("You must provide either of 'services' or 'task'")
	}
	return nil
}
func (c *ECSPortForwardOptions) Run(f factory.Factory, cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	awsConfig, err := f.BuildAWSConfig(ctx)
	if err != nil {
		log.WithError(err).Errorf("failed to load aws config")
		return err
	}

	ecsService := f.BuildECSService(awsConfig)
	sessionManager := f.BuildSessionManager(awsConfig)

	log.Infof("connecting to clusrter:%s task:%s container:%s", c.cluster, c.task, c.container)
	targetId, err := ecsService.GetTargetIDByTaskID(ctx, c.cluster, c.task, c.container)
	if err != nil {
		log.WithError(err).Errorf("failed to connect to container")
		return err
	}

	locals := strings.Split(c.local, ":")
	if len(locals) < 2 {
		return validation.NewInvalidInputError("local must follow the format of LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT")
	}

	region := cmd.Flags().Lookup(flags.RegionFlag).Value.String()
	if len(locals) == 2 {
		return sessionManager.PortForwardingSession(&sessions.PortForwardingInput{
			Region:     region,
			Target:     targetId,
			LocalPort:  locals[0],
			RemotePort: locals[1],
		})
	}

	return sessionManager.PortForwardingToRemoteHostSession(&sessions.PortForwardingToRemoteInput{
		Region:     region,
		Target:     targetId,
		LocalPort:  locals[0],
		RemoteHost: locals[1],
		RemotePort: locals[2],
	})
}

func New(f factory.Factory) *cobra.Command {
	op := &ECSPortForwardOptions{}
	command := &cobra.Command{
		Use:     "ecs-port-forward",
		Short:   synopsis,
		PreRun:  prerun.Setup,
		Example: "aws-ssh ecs-port-forward --cluster CLUSTER_NAME --local LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT --task TASK_ID",
		Long:    help,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.Init(
				cmd.Flags().Lookup(flags.AWSProfileFlag).Value.String(),
				cmd.Flags().Lookup(flags.RegionFlag).Value.String(),
			)
			if err := op.Validate(f, cmd, args); err != nil {
				return err
			}
			if err := op.Run(f, cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	command.Flags().StringVarP(
		&op.cluster,
		clusterFlag,
		"c",
		"",
		"ECS Cluster Name",
	)
	command.Flags().StringVarP(
		&op.service,
		serviceFlag,
		"s",
		"",
		"ECS Service Name. "+
			"If provided, it will search the ECS Service and try to access to an Active task. "+
			"Either of Service or Task must be provided.",
	)
	command.Flags().StringVarP(
		&op.task,
		taskFlag,
		"t",
		"",
		"ECS Task ID. "+
			"Either of Service or Task must be provided.",
	)
	command.Flags().StringVarP(
		&op.container,
		containerFlag,
		"n",
		"",
		"Container name. Required if task is running more than one container",
	)
	command.Flags().StringVarP(
		&op.local,
		localFlag,
		"L",
		"",
		"LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT "+
			"Forward a local port to a remote address/port",
	)

	command.MarkFlagsRequiredTogether(clusterFlag, localFlag)
	return command
}

const synopsis = "Port forwarding for AWS ECS tasks."
const help = `
Usage: aws-ssh ecs-port-forward [ECSPortForwardOptions]
	Forward localPort port to localPort port on task
	Forward localPort port to a remote host/port accessible from task
 	`
