package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Task struct {
	Arn        string
	Containers []Container
}

func NewTask(t awsTypes.Task) Task {

	var containers []Container
	for _, c := range t.Containers {
		containers = append(containers, NewContainer(c))
	}
	return Task{
		Arn:        aws.ToString(t.TaskArn),
		Containers: containers,
	}
}
