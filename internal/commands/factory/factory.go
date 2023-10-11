//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=factory.go -package=mock -destination=./mock/factory_mock.go
package factory

import (
	"aws-ssh/internal/ecs_ssh"
	"aws-ssh/internal/sessions"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Factory interface {
	Init(profile, region string)
	BuildAWSConfig(ctx context.Context) (aws.Config, error)
	BuildECSService(cfg aws.Config) ecs_ssh.ECSService
	BuildSessionManager(cfg aws.Config) sessions.Manager
}
