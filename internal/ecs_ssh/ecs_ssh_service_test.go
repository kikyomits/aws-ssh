package ecs_ssh_test

import (
	"aws-ssh/internal/ecs_ssh"
	"aws-ssh/internal/ecs_ssh/mock"
	"aws-ssh/internal/ecs_ssh/types"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestECSSSHService_GetTargetIDByTaskID(t *testing.T) {
	clusterName := "a-cluster"
	taskID := "a-task-id"
	containerName := "a-container"
	runtimeID := "a-runtime-id"
	targetID := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)

	type args struct {
		containerName string
	}
	tests := []struct {
		name    string
		args    args
		out     func(a args) (types.Task, error)
		want    string
		wantErr bool
	}{
		{
			name: "GIVEN correct containerName specified WHEN multiple containers are found THEN return TaskID",
			args: args{
				containerName: containerName,
			},
			out: func(a args) (types.Task, error) {
				return types.Task{
					Containers: []types.Container{
						{
							Name:      containerName,
							RuntimeID: runtimeID,
						},
						{
							Name:      "other container",
							RuntimeID: "other-runtime-id",
						},
					},
				}, nil
			},
			want: targetID,
		},
		{
			name: "GIVEN containerName not specified WHEN single container is found THEN return TaskID",
			args: args{},
			out: func(a args) (types.Task, error) {
				return types.Task{
					Containers: []types.Container{
						{
							Name:      "other container",
							RuntimeID: runtimeID,
						},
					},
				}, nil
			},
			want: targetID,
		},
		{
			name: "GIVEN containerName not specified WHEN multiple containers are found THEN return err",
			args: args{},
			out: func(a args) (types.Task, error) {
				return types.Task{
					Containers: []types.Container{
						{
							Name:      "other container",
							RuntimeID: runtimeID,
						},
						{
							Name:      containerName,
							RuntimeID: runtimeID,
						},
					},
				}, nil
			},
			wantErr: true,
		},
		{
			name: "WHEN AWS returns error not specified THEN return err",
			args: args{},
			out: func(a args) (types.Task, error) {
				return types.Task{}, fmt.Errorf("a-err")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			a := mock.NewMockECSAgent(ctrl)
			c := ecs_ssh.NewECSService(a)
			a.EXPECT().GetTask(gomock.Any(), ecs_ssh.GetTaskInput{TaskID: taskID, ClusterName: clusterName}).Return(tt.out(tt.args))
			got, err := c.GetTargetIDByTaskID(context.Background(), clusterName, taskID, tt.args.containerName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestECSSSHService_GetTargetIDByServiceName(t *testing.T) {
	clusterName := "a-cluster"
	serviceName := "a-service"
	taskID := "a-task-id"
	containerName := "a-container"
	runtimeID := "a-runtime-id"
	targetID := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)
	task := types.Task{
		Containers: []types.Container{
			{
				Name:      containerName,
				RuntimeID: runtimeID,
			},
		},
	}

	tests := []struct {
		name         string
		mock         func(m *mock.MockECSAgent)
		listTasksOut func() ([]string, error)
		getTaskOut   func() (types.Task, error)
		want         string
		wantErr      bool
	}{
		{
			name: "WHEN successfully list tasks and get task details THEN Should return target id",
			mock: func(m *mock.MockECSAgent) {
				m.EXPECT().ListRunningTask(gomock.Any(), ecs_ssh.ListRunningTasksInput{ServiceName: serviceName, ClusterName: clusterName}).Return([]string{taskID}, nil)
				m.EXPECT().GetTask(gomock.Any(), ecs_ssh.GetTaskInput{TaskID: taskID, ClusterName: clusterName}).Return(task, nil)
			},
			want: targetID,
		},
		{
			name: "WHEN failed to list tasks and get task details THEN Should return err",
			mock: func(m *mock.MockECSAgent) {
				m.EXPECT().ListRunningTask(gomock.Any(), ecs_ssh.ListRunningTasksInput{ServiceName: serviceName, ClusterName: clusterName}).Return([]string{}, fmt.Errorf("a-err"))
			},
			wantErr: true,
		},
		{
			name: "WHEN failed to get task details THEN Should return err",
			mock: func(m *mock.MockECSAgent) {
				m.EXPECT().ListRunningTask(gomock.Any(), ecs_ssh.ListRunningTasksInput{ServiceName: serviceName, ClusterName: clusterName}).Return([]string{taskID}, nil)
				m.EXPECT().GetTask(gomock.Any(), ecs_ssh.GetTaskInput{TaskID: taskID, ClusterName: clusterName}).Return(types.Task{}, fmt.Errorf("a-err"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			a := mock.NewMockECSAgent(ctrl)
			c := ecs_ssh.NewECSService(a)
			tt.mock(a)

			got, err := c.GetTargetIDByServiceName(context.Background(), clusterName, serviceName, containerName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
