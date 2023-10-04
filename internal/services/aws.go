package services

import (
	"aws-ssh/internal/agents"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithHTTPClient(agents.NewHTTPClient()),
	)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}
