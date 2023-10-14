package ecs_exec_test

import (
	"aws-ssh/internal/commands"
	"aws-ssh/internal/commands/ecs_exec"
	"aws-ssh/internal/commands/factory/mock"
	mock3 "aws-ssh/internal/ecs_ssh/mock"
	"aws-ssh/internal/sessions"
	mock2 "aws-ssh/internal/sessions/mock"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestECSPortForwardOptions_Run(t *testing.T) {
	cluster := "a-cluster"
	service := "a-service"
	task := "a-task"
	container := "a-container"
	cmds := []string{"echo", "abc"}
	command := "echo abc"
	targetID := "ecs:a_b_c"

	type args struct {
		service string
		task    string
	}

	tests := []struct {
		name    string
		mock    func(f *mock.MockFactory, sm *mock2.MockManager, ecs *mock3.MockECSService)
		args    args
		wantErr bool
	}{
		{
			name: "GIVEN task ID is provided THEN find target by task ID and execute command",
			args: args{
				task: task,
			},
			mock: func(f *mock.MockFactory, sm *mock2.MockManager, ecs *mock3.MockECSService) {
				f.EXPECT().BuildAWSConfig(gomock.Any()).Return(aws.Config{}, nil)
				f.EXPECT().BuildSessionManager(gomock.Any()).Return(sm)
				f.EXPECT().BuildECSService(gomock.Any()).Return(ecs)
				ecs.EXPECT().GetTargetIDByTaskID(gomock.Any(), cluster, task, container).Return(targetID, nil)
				in := sessions.NewExecInput("", targetID, command)
				sm.EXPECT().ExecSession(in).Return(nil)
			},
		},
		{
			name: "GIVEN service is provided THEN find target by task ID and execute command",
			args: args{
				service: service,
			},
			mock: func(f *mock.MockFactory, sm *mock2.MockManager, ecs *mock3.MockECSService) {
				f.EXPECT().BuildAWSConfig(gomock.Any()).Return(aws.Config{}, nil)
				f.EXPECT().BuildSessionManager(gomock.Any()).Return(sm)
				f.EXPECT().BuildECSService(gomock.Any()).Return(ecs)
				ecs.EXPECT().GetTargetIDByServiceName(gomock.Any(), cluster, service, container).Return(targetID, nil)
				in := sessions.NewExecInput("", targetID, command)
				sm.EXPECT().ExecSession(in).Return(nil)
			},
		},
		{
			name: "GIVEN invalid aws config WHEN BuildAWSConfig THEN return err",
			args: args{
				service: service,
			},
			mock: func(f *mock.MockFactory, sm *mock2.MockManager, ecs *mock3.MockECSService) {
				f.EXPECT().BuildAWSConfig(gomock.Any()).Return(aws.Config{}, fmt.Errorf("a-err"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sm := mock2.NewMockManager(ctrl)
			ecs := mock3.NewMockECSService(ctrl)
			f := mock.NewMockFactory(ctrl)
			tt.mock(f, sm, ecs)
			cmd, err := commands.RegisteredCommands()
			assert.NoError(t, err)

			c := &ecs_exec.ECSExecOptions{
				Cluster:   cluster,
				Service:   tt.args.service,
				Task:      tt.args.task,
				Container: container,
			}

			err = c.Run(f, &cmd, cmds)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestECSPortForwardOptions_Validate(t *testing.T) {
	cluster := "a-cluster"
	service := "a-service"
	task := "a-task"
	type fields struct {
		Service   string
		Task      string
		Container string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "GIVEN service name provided WHEN port forward THEN return nil",
			fields: fields{
				Service: service,
			},
		},
		{
			name: "GIVEN task name provided WHEN port forward THEN return nil",
			fields: fields{
				Task: task,
			},
		},
		{
			name:    "GIVEN task and service names are missing WHEN port forward THEN return error",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := mock.NewMockFactory(ctrl)
			c := &ecs_exec.ECSExecOptions{
				Cluster:   cluster,
				Service:   tt.fields.Service,
				Task:      tt.fields.Task,
				Container: tt.fields.Container,
			}
			err := c.Validate(f, &cobra.Command{}, []string{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
