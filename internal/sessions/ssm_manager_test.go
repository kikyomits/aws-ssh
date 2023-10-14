package sessions_test

import (
	"aws-ssh/internal/sessions"
	"aws-ssh/internal/sessions/mock"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"go.uber.org/mock/gomock"
)

func TestSSMSessionManager_PortForwardingSession(t *testing.T) {
	opts := sessions.NewPortForwardingInput("ap-southeast-2", "a-target", "8080", "8080")

	type args struct {
		startSessionOutput func() (*ssm.StartSessionOutput, error)
		pluginOutput       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GIVEN successful connection WHEN start session THEN return nil",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return &ssm.StartSessionOutput{
						SessionId:  aws.String("a-session-id"),
						StreamUrl:  aws.String("a-stream-url"),
						TokenValue: aws.String("a-token"),
					}, nil
				},
				pluginOutput: nil,
			},
		},
		{
			name: "GIVEN invalid session WHEN start session THEN return err",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return nil, errors.New("err")
				},
			},
			wantErr: true,
		},
		{
			name: "GIVEN plugin failure WHEN start session THEN return err",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return &ssm.StartSessionOutput{
						SessionId:  aws.String("a-session-id"),
						StreamUrl:  aws.String("a-stream-url"),
						TokenValue: aws.String("a-token"),
					}, nil
				},
				pluginOutput: errors.New("err"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ssmClient := mock.NewMockSSM(ctrl)
			plugin := mock.NewMockPlugin(ctrl)

			ssmClient.EXPECT().StartSession(gomock.Any(), gomock.Any()).Return(tt.args.startSessionOutput())
			plugin.EXPECT().Execute(gomock.Any()).AnyTimes().Return(tt.args.pluginOutput)
			c := sessions.NewSSMSessionManager(plugin, ssmClient)
			if err := c.PortForwardingSession(opts); (err != nil) != tt.wantErr {
				t.Errorf("PortForwardingSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSSMSessionManager_PortForwardingToRemoteHostSession(t *testing.T) {

	opts := sessions.NewPortForwardingToRemoteInput("ap-southeast-2", "a-target", "8080", "8080", "https://example.com")

	type args struct {
		startSessionOutput func() (*ssm.StartSessionOutput, error)
		pluginOutput       error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GIVEN successful connection WHEN start session THEN return nil",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return &ssm.StartSessionOutput{
						SessionId:  aws.String("a-session-id"),
						StreamUrl:  aws.String("a-stream-url"),
						TokenValue: aws.String("a-token"),
					}, nil
				},
				pluginOutput: nil,
			},
		},
		{
			name: "GIVEN invalid session WHEN start session THEN return err",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return nil, errors.New("err")
				},
			},
			wantErr: true,
		},
		{
			name: "GIVEN plugin failure WHEN start session THEN return err",
			args: args{
				startSessionOutput: func() (*ssm.StartSessionOutput, error) {
					return &ssm.StartSessionOutput{
						SessionId:  aws.String("a-session-id"),
						StreamUrl:  aws.String("a-stream-url"),
						TokenValue: aws.String("a-token"),
					}, nil
				},
				pluginOutput: errors.New("err"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		ssmClient := mock.NewMockSSM(ctrl)
		plugin := mock.NewMockPlugin(ctrl)

		t.Run(tt.name, func(t *testing.T) {
			ssmClient.EXPECT().StartSession(gomock.Any(), gomock.Any()).Return(tt.args.startSessionOutput())
			plugin.EXPECT().Execute(gomock.Any()).AnyTimes().Return(tt.args.pluginOutput)
			c := sessions.NewSSMSessionManager(plugin, ssmClient)
			if err := c.PortForwardingToRemoteHostSession(opts); (err != nil) != tt.wantErr {
				t.Errorf("PortForwardingSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
