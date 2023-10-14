package ecs_ssh_test

import (
	"aws-ssh/internal/ecs_ssh"
	"aws-ssh/internal/ecs_ssh/mock"
	"aws-ssh/internal/ecs_ssh/types"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestECSSSHAgent_GetTask(t *testing.T) {

	clusterName := "a-cluster"
	taskID := "a-task"
	containerName := "a-container"
	runtimeID := "a_runtime_id"
	type args struct {
		containerName       string
		describeTasksOutput func() (*ecs.DescribeTasksOutput, error)
	}
	tests := []struct {
		name    string
		arg     args
		want    types.Task
		wantErr bool
	}{
		{
			name: "GIVEN single container task WHEN find task THEN return runtime id",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []awsTypes.Task{
							{
								Containers: []awsTypes.Container{
									{
										Name:      aws.String(containerName),
										RuntimeId: aws.String(runtimeID),
									},
								},
							},
						},
					}, nil
				},
			},
			want: types.Task{
				Arn: "",
				Containers: []types.Container{
					{
						Name:      containerName,
						RuntimeID: runtimeID,
					},
				},
			},
		},
		{
			name: "WHEN no task found THEN return error",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []awsTypes.Task{},
					}, nil
				},
				containerName: containerName,
			},
			wantErr: true,
		},
		{
			name: "WHEN error thrown THEN return error",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{}, fmt.Errorf("a error")
				},
				containerName: containerName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := mock.NewMockECS(ctrl)
			c.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(tt.arg.describeTasksOutput())
			a := &ecs_ssh.ECSSSHClient{
				ECS: c,
			}
			got, err := a.GetTask(context.Background(), ecs_ssh.GetTaskInput{TaskID: taskID, ClusterName: clusterName})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetIDByTaskID() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestECSSSHAgent_ListRunningTask(t *testing.T) {

	clusterName := "a-cluster"
	containerName := "a-container"
	taskArns := []string{"a-task", "b-task"}
	type args struct {
		containerName   string
		listTasksOutput func() (*ecs.ListTasksOutput, error)
	}
	tests := []struct {
		name    string
		arg     args
		want    []string
		wantErr bool
	}{
		{
			name: "GIVEN single container task WHEN find task THEN return runtime id",
			arg: args{
				listTasksOutput: func() (*ecs.ListTasksOutput, error) {
					return &ecs.ListTasksOutput{
						TaskArns: taskArns,
					}, nil
				},
			},
			want: taskArns,
		},
		{
			name: "WHEN no task found THEN return error",
			arg: args{
				listTasksOutput: func() (*ecs.ListTasksOutput, error) {
					return &ecs.ListTasksOutput{
						TaskArns: []string{},
					}, nil
				},
				containerName: containerName,
			},
			wantErr: true,
		},
		{
			name: "WHEN error thrown THEN return error",
			arg: args{
				listTasksOutput: func() (*ecs.ListTasksOutput, error) {
					return &ecs.ListTasksOutput{
						TaskArns: []string{"a"},
					}, fmt.Errorf("a error")
				},
				containerName: containerName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := mock.NewMockECS(ctrl)
			c.EXPECT().ListTasks(gomock.Any(), gomock.Any()).Return(tt.arg.listTasksOutput())
			a := &ecs_ssh.ECSSSHClient{
				ECS: c,
			}
			got, err := a.ListRunningTasks(context.Background(), ecs_ssh.ListRunningTasksInput{ClusterName: clusterName, ServiceName: "a-service"})
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
