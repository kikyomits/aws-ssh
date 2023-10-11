package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Container struct {
	Name      string
	RuntimeID string
}

func NewContainer(c awsTypes.Container) Container {
	return Container{
		Name:      aws.ToString(c.Name),
		RuntimeID: aws.ToString(c.RuntimeId),
	}
}
