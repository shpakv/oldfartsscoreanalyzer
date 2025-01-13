package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var cachedDynamoDB *dynamodb.Client

func DynamoDB(ctx context.Context) *dynamodb.Client {
	if cachedDynamoDB == nil {
		cfg := mustLoadConfig(ctx)
		cachedDynamoDB = dynamodb.NewFromConfig(cfg)
	}
	return cachedDynamoDB
}
