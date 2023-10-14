package ecs_exec

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
	clusterFlag   = "cluster"
	serviceFlag   = "service"
	taskFlag      = "task"
	containerFlag = "container"
	commandFlag   = "command"
)

type ECSExecOptions struct {
	Cluster   string
	Service   string
	Task      string
	Container string
}

func (c *ECSExecOptions) Validate(_ factory.Factory, _ *cobra.Command, _ []string) error {
	if c.Service == "" && c.Task == "" {
		return validation.NewInvalidInputError("You must provide either of 'service' or 'task'")
	}

	if c.Service != "" && c.Task != "" {
		return validation.NewInvalidInputError("You cannot specify service name and task ID at the same time. Please specify either of them")
	}

	return nil
}
func (c *ECSExecOptions) Run(f factory.Factory, cmd *cobra.Command, args []string) error {
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

	region, err := cmd.Flags().GetString(flags.RegionFlag)
	if err != nil {
		log.WithError(err).Debugf("failed to get region option. ignoring this erro")
		region = ""
	}
	return sessionManager.ExecSession(sessions.NewExecInput(region, targetID, strings.Join(args, " ")))
}

func (c *ECSExecOptions) getTargetID(ctx context.Context, f factory.Factory, awsConfig aws.Config) (string, error) {
	ecsService := f.BuildECSService(awsConfig)
	if c.Service != "" {
		return ecsService.GetTargetIDByServiceName(ctx, c.Cluster, c.Service, c.Container)
	}
	return ecsService.GetTargetIDByTaskID(ctx, c.Cluster, c.Task, c.Container)
}

func New(f factory.Factory) *cobra.Command {
	op := &ECSExecOptions{}
	command := &cobra.Command{
		Use:     "ecs-exec [flags] COMMAND",
		Short:   synopsis,
		PreRun:  prerun.Setup,
		Example: "aws-ssh ecs-exec --cluster CLUSTER_NAME --task TASK_ID bash",
		Args:    cobra.MinimumNArgs(1),
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
				log.WithError(err).Errorf("command failed")
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

	command.MarkFlagsRequiredTogether(clusterFlag)
	return command
}

const synopsis = "Execute interactive command in AWS ECS task."
