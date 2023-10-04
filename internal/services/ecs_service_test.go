package services

import (
	"aws-ssh/internal/services/mock"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/mock/gomock"
)

func TestECSService_GetECSRuntimeIDByTaskID(t *testing.T) {

	ctrl := gomock.NewController(t)
	ecsClient := mock.NewMockECSClient(ctrl)
	clusterName := "a-cluster"
	taskID := "a-task"
	containerName := "a-container"
	runtimeID := "a_runtime_id"
	targetID := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, runtimeID)
	type args struct {
		containerName       string
		describeTasksOutput func() (*ecs.DescribeTasksOutput, error)
	}
	tests := []struct {
		name    string
		arg     args
		want    string
		wantErr bool
	}{
		{
			name: "GIVEN single container task WHEN find task THEN return runtime id",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []types.Task{
							{
								Containers: []types.Container{
									{
										Name:      aws.String("other-container"),
										RuntimeId: aws.String(runtimeID),
									},
								},
							},
						},
					}, nil
				},
			},
			want: targetID,
		},
		{
			name: "GIVEN multiple container task WHEN find task THEN filter by container name and return runtime id",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []types.Task{
							{
								Containers: []types.Container{
									{
										Name:      aws.String("other-container"),
										RuntimeId: aws.String("other-runtime-id"),
									},
									{
										Name:      aws.String(containerName),
										RuntimeId: aws.String(runtimeID),
									},
								},
							},
						},
					}, nil
				},
				containerName: containerName,
			},
			want: targetID,
		},
		{
			name: "WHEN no task found THEN return error",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []types.Task{},
					}, nil
				},
				containerName: containerName,
			},
			wantErr: true,
		},
		{
			name: "WHEN no target container found THEN return error",
			arg: args{
				describeTasksOutput: func() (*ecs.DescribeTasksOutput, error) {
					return &ecs.DescribeTasksOutput{
						Tasks: []types.Task{
							{
								Containers: []types.Container{
									{
										Name:      aws.String("other-container"),
										RuntimeId: aws.String("other-runtime-id"),
									},
									{
										Name:      aws.String("another-container"),
										RuntimeId: aws.String("another-runtime-id"),
									},
								},
							},
						},
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
			ecsClient.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Return(tt.arg.describeTasksOutput())
			c := &ECSSSHService{
				ecs: ecsClient,
			}
			got, err := c.GetTargetIDByTaskID(context.Background(), clusterName, taskID, tt.arg.containerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetIDByTaskID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTargetIDByTaskID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
