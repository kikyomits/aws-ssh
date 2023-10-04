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

type ECSService interface {
	GetTargetIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error)
}
type ECSSSHService struct {
	ecs ECSClient
}

func NewECSService(ecsClient ECSClient) *ECSSSHService {
	return &ECSSSHService{
		ecs: ecsClient,
	}
}

func (c *ECSSSHService) findTaskIDByServiceName(ctx context.Context, serviceName string) (*string, error) {
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
		return nil, fmt.Errorf(
			"found multiple ECS Services, please provide unique name. found: %v", names)
	}

	log.Debugf("found ECS Service %s", serviceName)
	return out.Services[0].TaskSets[0].TaskSetArn, nil

}

func (c *ECSSSHService) findContainerByName(containers []types.Container, name string) (*types.Container, error) {
	if len(containers) > 1 && name == "" {
		var names []*string
		for _, container := range containers {
			names = append(names, container.Name)
		}
		return nil, fmt.Errorf("Need to specify the container name. found: %v", names)
	}

	if len(containers) == 1 && name == "" {
		return &containers[0], nil
	}

	for _, container := range containers {
		if *container.Name == name {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("cannot find container '%s' in the given task", name)
}

func (c *ECSSSHService) GetTargetIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error) {
	input := &ecs.DescribeTasksInput{
		Tasks:   []string{taskID},
		Cluster: aws.String(clusterName),
	}
	tasks, err := c.ecs.DescribeTasks(ctx, input)
	if err != nil {
		return "", err
	}

	var taskArns []string
	for _, t := range tasks.Tasks {
		taskArns = append(taskArns, aws.ToString(t.TaskArn))
	}

	log.WithField("tasks", taskArns).Debugf("found tasks")
	if len(tasks.Tasks) == 0 {
		return "", fmt.Errorf("cannot find any tasks by taskID: %s", taskID)
	}

	task := tasks.Tasks[0]
	containers := task.Containers
	if len(containers) == 0 {
		return "", fmt.Errorf("cannot find any containers in taskID: %s", taskID)
	}

	container, err := c.findContainerByName(containers, containerName)
	if err != nil {
		return "", err
	}
	log.Debugf("found container %s", containerName)
	runtimeID := aws.ToString(container.RuntimeId)
	targetID := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)
	log.Debugf("found container runtime id:%s", runtimeID)
	return targetID, nil
}
