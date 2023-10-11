//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=ecs_service.go  -package=mock  -destination=./mock/ecs_service_mock.go
package ecs_ssh

import (
	"context"
)

type ECSService interface {
	GetTargetIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error)
	GetTargetIDByServiceName(ctx context.Context, clusterName, serviceName, containerName string) (string, error)
}
