package ecs_port_forward

import (
	"aws-ssh/internal/commands/flags"
	"aws-ssh/internal/commands/prerun"
	"aws-ssh/internal/services"
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

type options struct {
	cluster   string
	service   string
	task      string
	container string
	local     string
}

func New() *cobra.Command {
	op := &options{}
	command := &cobra.Command{
		Use:    "ecs_port_forward",
		Short:  synopsis,
		PreRun: prerun.Setup,
		Long:   help,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			if op.service == "" && op.task == "" {
				return validation.NewInvalidInputError("You must provide either of 'services' or 'task'")
			}

			region := cmd.Flags().Lookup(flags.RegionFlag).Value.String()
			log.Infof("connecting to clusrter %s task %s container %s", op.cluster, op.task, op.container)

			cfg, err := services.NewAWSConfig(ctx, region)
			if err != nil {
				return err
			}
			e, err := services.NewECSService(cfg)
			if err != nil {
				return err
			}

			targetId, err := e.GetECSRuntimeIDByTaskID(ctx, op.cluster, op.task, op.container)
			if err != nil {
				return err
			}
			c, err := sessions.NewSSMSessionManager(cfg)
			if err != nil {
				return err
			}

			locals := strings.Split(op.local, ":")
			if len(locals) < 2 {
				return validation.NewInvalidInputError("local must follow the format of LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT")
			} else if len(locals) == 2 {
				err = c.PortForwardingSession(&sessions.PortForwardingInput{
					Target:     targetId,
					LocalPort:  locals[0],
					RemotePort: locals[1],
				})
				if err != nil {
					return err
				}
			} else {
				err = c.PortForwardingToRemoteHostSession(&sessions.PortForwardingToRemoteInput{
					Target:     targetId,
					LocalPort:  locals[0],
					Host:       locals[1],
					RemotePort: locals[2],
				})
				if err != nil {
					return err
				}
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
Usage: aws-ssh ecs_port_forward [options]
	Forward localPort port to localPort port on task
	Forward localPort port to a remote host/port accessible from task
 	`
