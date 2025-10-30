package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cachedConfig *aws.Config

func mustLoadConfig(ctx context.Context) aws.Config {
	if cachedConfig == nil {
		newConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic(fmt.Sprintf("could not load AWS configuration: %v", err))
		}
		cachedConfig = &newConfig
	}
	return *cachedConfig
}
