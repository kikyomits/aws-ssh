//go:generate go run go.uber.org/mock/mockgen@v0.3.0 -source=ssm.go -package=mock -destination=./mock/ssm_mock.go
package sessions

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSM interface {
	StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)
}
