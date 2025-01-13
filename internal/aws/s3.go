package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var cachedS3 *s3.Client

func S3(ctx context.Context) *s3.Client {
	if cachedS3 == nil {
		cfg := mustLoadConfig(ctx)
		cachedS3 = s3.NewFromConfig(cfg)
	}
	return cachedS3
}
