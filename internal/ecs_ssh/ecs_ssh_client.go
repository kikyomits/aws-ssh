package ecs_ssh

import (
	"aws-ssh/internal/ecs_ssh/types"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSSSHClient struct {
	ECS ECS
}

func NewECSSSHAgent(client ECS) ECSSSHClient {
	return ECSSSHClient{ECS: client}
}

func (a *ECSSSHClient) GetTask(ctx context.Context, in GetTaskInput) (types.Task, error) {
	input := &ecs.DescribeTasksInput{
		Tasks:   []string{in.TaskID},
		Cluster: aws.String(in.ClusterName),
	}
	tasks, err := a.ECS.DescribeTasks(ctx, input)
	if err != nil {
		return types.Task{}, err
	}

	if len(tasks.Tasks) == 0 {
		return types.Task{}, fmt.Errorf("cannot find any tasks by taskID: %s", in.TaskID)
	}

	return types.NewTask(tasks.Tasks[0]), nil
}

func (a *ECSSSHClient) ListRunningTasks(ctx context.Context, in ListRunningTasksInput) ([]string, error) {
	req := &ecs.ListTasksInput{
		Cluster:       aws.String(in.ClusterName),
		ServiceName:   aws.String(in.ServiceName),
		DesiredStatus: awsTypes.DesiredStatusRunning,
	}
	res, err := a.ECS.ListTasks(ctx, req)
	if err != nil {
		return []string{}, err
	}

	if len(res.TaskArns) == 0 {
		return []string{}, fmt.Errorf("cannot find any RUNNING tasks in service: %s", in.ServiceName)
	}
	return res.TaskArns, nil
}
