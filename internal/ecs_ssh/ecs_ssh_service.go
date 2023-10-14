package ecs_ssh

import (
	"aws-ssh/internal/ecs_ssh/types"
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	SSMDocumentAWSStartPortForwardingSessionToRemoteHost = "AWS-StartPortForwardingSessionToRemoteHost"
	SSMDocumentAWSStartPortForwardingSession             = "AWS-StartPortForwardingSession"
	SSMStartInteractiveCommand                           = "AWS-StartInteractiveCommand"
)

type ECSSSHService struct {
	ecs ECSClient
}

func NewECSService(ecsAgent ECSClient) *ECSSSHService {
	return &ECSSSHService{
		ecs: ecsAgent,
	}
}

func (c *ECSSSHService) GetTargetIDByServiceName(ctx context.Context, clusterName, serviceName, containerName string) (string, error) {
	log.WithField("service", serviceName).Debugf("finding target Task by Service name")
	arns, err := c.ecs.ListRunningTasks(ctx, ListRunningTasksInput{ClusterName: clusterName, ServiceName: serviceName})
	if err != nil {
		return "", errors.Wrap(err, "failed to find a Task by Service Name")
	}

	arn := arns[0]
	if len(arns) > 1 {
		log.Infof("multiple RUNNING tasks are found. automatically select %s", arn)
	}

	log.WithField("task", arn).Debugf("found RUNNING task")
	sp := strings.Split(arn, "/")
	taskID := sp[len(sp)-1]

	log.Infof("connecting to taskID %s", taskID)
	return c.GetTargetIDByTaskID(ctx, clusterName, taskID, containerName)
}

func (c *ECSSSHService) GetTargetIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error) {
	log.WithField("task", taskID).Debugf("finding target Task by Task ID")

	task, err := c.ecs.GetTask(ctx, GetTaskInput{ClusterName: clusterName, TaskID: taskID})
	if err != nil {
		return "", errors.Wrap(err, "failed to identify a Task")
	}

	containers := task.Containers

	container, err := c.findContainerByName(containers, containerName)
	if err != nil {
		return "", errors.Wrap(err, "failed to identify a Task")
	}

	runtimeID := container.RuntimeID
	targetID := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)
	log.WithFields(map[string]interface{}{
		"runtimeID":     runtimeID,
		"containerName": containerName,
	}).Debugf("found container")
	return targetID, nil
}

func (c *ECSSSHService) findContainerByName(containers []types.Container, name string) (*types.Container, error) {
	if len(containers) == 0 {
		return nil, fmt.Errorf("cannot find any containers in the specified task")
	}

	if len(containers) > 1 && name == "" {
		var names []string
		for _, container := range containers {
			names = append(names, container.Name)
		}
		return nil, fmt.Errorf("Need to specify the container name. found: %v", names)
	}

	if len(containers) == 1 && name == "" {
		return &containers[0], nil
	}

	for _, container := range containers {
		if container.Name == name {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("cannot find container '%s' in the given task", name)
}
