//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=ecs_agent.go  -package=mock  -destination=./mock/ecs_agent_mock.go
package ecs_ssh

import (
	"aws-ssh/internal/ecs_ssh/types"
	"context"
)

type ECSAgent interface {
	ListRunningTasks(ctx context.Context, in ListRunningTasksInput) ([]string, error)
	GetTask(ctx context.Context, in GetTaskInput) (types.Task, error)
}

type GetTaskInput struct {
	ClusterName string
	TaskID      string
}

type ListRunningTasksInput struct {
	ClusterName string
	ServiceName string
}
