package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/sirupsen/logrus"
)

const (
	SSMDocumentAWSStartPortForwardingSessionToRemoteHost = "AWS-StartPortForwardingSessionToRemoteHost"
	SSMDocumentAWSStartPortForwardingSession             = "AWS-StartPortForwardingSession"
)

type ECSService struct {
	ecs ECSClient
}

func NewECSService(ecsClient ECSClient) (*ECSService, error) {
	return &ECSService{
		ecs: ecsClient,
	}, nil
}

func (c *ECSService) findTaskIDByServiceName(ctx context.Context, serviceName string) (*string, error) {
	input := &ecs.DescribeServicesInput{Services: []string{serviceName}}
	out, err := c.ecs.DescribeServices(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(out.Services) == 0 {
		return nil, errors.New("cannot find any services from given name")
	}

	if len(out.Services) > 1 {
		var names []*string
		for _, svc := range out.Services {
			names = append(names, svc.ServiceName)
		}
		return nil, errors.New(fmt.Sprintf(
			"found multiple ECS Services, please provide unique name. found: %v", names))
	}

	log.Debugf("found ECS Service %s", serviceName)
	return out.Services[0].TaskSets[0].TaskSetArn, nil

}

func (c *ECSService) findContainerByName(containers []types.Container, name string) (*types.Container, error) {
	if len(containers) > 1 && name == "" {
		var names []*string
		for _, container := range containers {
			names = append(names, container.Name)
		}
		return nil, errors.New(fmt.Sprintf("Need to specify the container name. found: %v", names))
	}

	if len(containers) == 0 && name == "" {
		return &containers[0], nil
	}

	for _, container := range containers {
		if *container.Name == name {
			return &container, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("cannot find container '%s' in the given task", name))
}

func (c *ECSService) GetECSRuntimeIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error) {
	input := &ecs.DescribeTasksInput{
		Tasks:   []string{taskID},
		Cluster: aws.String(clusterName),
	}
	tasks, err := c.ecs.DescribeTasks(ctx, input)
	if err != nil {
		return "", err
	}

	if len(tasks.Tasks) == 0 {
		return "", errors.New(fmt.Sprintf("cannot find any tasks by taskID: %s", taskID))
	}

	task := tasks.Tasks[0]
	containers := task.Containers

	container, err := c.findContainerByName(containers, containerName)

	if err != nil {
		return "", err
	}

	return aws.ToString(container.RuntimeId), nil
}
