package ecs_port_forward

import (
	"aws-ssh/internal/commands/factory"
	"aws-ssh/internal/commands/flags"
	"aws-ssh/internal/commands/prerun"
	"aws-ssh/internal/sessions"
	"aws-ssh/internal/validation"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	clusterFlag   = "Cluster"
	serviceFlag   = "ecs_ssh"
	taskFlag      = "Task"
	localFlag     = "Local"
	containerFlag = "Container"
)

type ECSPortForwardOptions struct {
	Cluster   string
	Service   string
	Task      string
	Container string
	Local     string
}

func (c *ECSPortForwardOptions) Validate(_ factory.Factory, _ *cobra.Command, _ []string) error {
	if c.Service == "" && c.Task == "" {
		return validation.NewInvalidInputError("You must provide either of 'Service' or 'Task'")
	}

	if c.Service != "" && c.Task != "" {
		return validation.NewInvalidInputError("You cannot specify service name and task ID at the same time. Please specify either of them")
	}

	locals := strings.Split(c.Local, ":")
	if len(locals) < 2 {
		return validation.NewInvalidInputError("Local must follow the format of LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT")
	}
	return nil
}
func (c *ECSPortForwardOptions) Run(f factory.Factory, cmd *cobra.Command, _ []string) error {
	ctx := context.Background()

	awsConfig, err := f.BuildAWSConfig(ctx)
	if err != nil {
		log.WithError(err).Errorf("failed to load aws config")
		return err
	}
	sessionManager := f.BuildSessionManager(awsConfig)
	targetID, err := c.getTargetID(ctx, f, awsConfig)
	if err != nil {
		return err
	}

	log.Infof("connecting to clusrter:%s Task:%s Container:%s", c.Cluster, c.Task, c.Container)

	// locals length must be validated in the Validate method
	locals := strings.Split(c.Local, ":")
	region, err := cmd.Flags().GetString(flags.RegionFlag)
	if err != nil {
		log.WithError(err).Debugf("failed to get region option. ignoring this erro")
		region = ""
	}
	if len(locals) == 2 {
		return sessionManager.PortForwardingSession(&sessions.PortForwardingInput{
			Region:     region,
			Target:     targetID,
			LocalPort:  locals[0],
			RemotePort: locals[1],
		})
	}

	return sessionManager.PortForwardingToRemoteHostSession(&sessions.PortForwardingToRemoteInput{
		Region:     region,
		Target:     targetID,
		LocalPort:  locals[0],
		RemoteHost: locals[1],
		RemotePort: locals[2],
	})
}

func (c *ECSPortForwardOptions) getTargetID(ctx context.Context, f factory.Factory, awsConfig aws.Config) (string, error) {
	ecsService := f.BuildECSService(awsConfig)
	if c.Service != "" {
		log.WithField("Service", c.Service).Infof("finding target Task by Service name")
		targetID, err := ecsService.GetTargetIDByServiceName(ctx, c.Cluster, c.Service, c.Container)
		if err != nil {
			log.WithError(err).Errorf("failed to find a Task by Service Name")
		}
		return targetID, err
	}

	targetID, err := ecsService.GetTargetIDByTaskID(ctx, c.Cluster, c.Task, c.Container)
	if err != nil {
		log.WithError(err).Errorf("failed to find a Task by TaskID")
		return "", err
	}
	return targetID, err
}

func New(f factory.Factory) *cobra.Command {
	op := &ECSPortForwardOptions{}
	command := &cobra.Command{
		Use:     "ecs-port-forward",
		Short:   synopsis,
		PreRun:  prerun.Setup,
		Example: "aws-ssh ecs-port-forward --Cluster CLUSTER_NAME --Local LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT --Task TASK_ID",
		Long:    help,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.Init(
				cmd.Flags().Lookup(flags.AWSProfileFlag).Value.String(),
				cmd.Flags().Lookup(flags.RegionFlag).Value.String(),
			)
			if err := op.Validate(f, cmd, args); err != nil {
				log.WithError(err).Errorf("validation failed")
				return err
			}
			if err := op.Run(f, cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	command.Flags().StringVarP(
		&op.Cluster,
		clusterFlag,
		"c",
		"",
		"ECS Cluster Name",
	)
	command.Flags().StringVarP(
		&op.Service,
		serviceFlag,
		"s",
		"",
		"ECS Service Name. "+
			"If provided, it will search the ECS Service and try to access to an Active Task. "+
			"Either of Service or Task must be provided.",
	)
	command.Flags().StringVarP(
		&op.Task,
		taskFlag,
		"t",
		"",
		"ECS Task ID. "+
			"Either of Service or Task must be provided.",
	)
	command.Flags().StringVarP(
		&op.Container,
		containerFlag,
		"n",
		"",
		"Container name. Required if Task is running more than one Container",
	)
	command.Flags().StringVarP(
		&op.Local,
		localFlag,
		"L",
		"",
		"LOCAL_PORT[:REMOTE_ADDR]:REMOTE_PORT "+
			"Forward a Local port to a remote address/port",
	)

	command.MarkFlagsRequiredTogether(clusterFlag, localFlag)
	return command
}

const synopsis = "Port forwarding for AWS ECS tasks."
const help = `
Usage: aws-ssh ecs_ssh-port-forward [ECSPortForwardOptions]
	Forward localPort port to localPort port on Task
	Forward localPort port to a remote host/port accessible from Task
 	`
