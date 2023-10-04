//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=ssm_client.go -package=mock -destination=./mock/ssm_client_mock.go
package sessions

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient interface {
	StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)
}
