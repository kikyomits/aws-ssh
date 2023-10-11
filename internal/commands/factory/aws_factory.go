package factory

import (
	"aws-ssh/internal/agents"
	"aws-ssh/internal/ecs_ssh"
	"aws-ssh/internal/sessions"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	log "github.com/sirupsen/logrus"
)

type AWSFactory struct {
	Profile string
	Region  string
}

func (f *AWSFactory) Init(profile, region string) {
	f.Profile = profile
	f.Region = region
}
func (f *AWSFactory) BuildAWSConfig(ctx context.Context) (aws.Config, error) {
	log.WithFields(map[string]interface{}{
		"region":  f.Region,
		"profile": f.Profile,
	}).Debug("load aws config")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(f.Region),
		config.WithSharedConfigProfile(f.Profile),
		config.WithHTTPClient(agents.NewHTTPClient()),
	)
	log.Debug(cfg.Region)
	log.Debug(cfg.AppID)

	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func (f *AWSFactory) BuildECSService(cfg aws.Config) ecs_ssh.ECSService {
	ecsClient := ecs.NewFromConfig(cfg)
	agent := ecs_ssh.NewECSSSHAgent(ecsClient)
	return ecs_ssh.NewECSService(&agent)
}

func (f *AWSFactory) BuildSessionManager(cfg aws.Config) sessions.Manager {
	return sessions.NewSSMSessionManager(sessions.NewAWSPluginSession(), ssm.NewFromConfig(cfg))
}
