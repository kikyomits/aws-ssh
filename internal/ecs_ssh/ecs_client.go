//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=ecs_client.go  -package=mock  -destination=./mock/ecs_client_mock.go
package ecs_ssh

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSClient interface {
	DescribeServices(ctx context.Context, params *ecs.DescribeServicesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)

	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
}
