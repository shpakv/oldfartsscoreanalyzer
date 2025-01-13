package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Tabler interface {
	TableName() string
}

type Repository[Entity interface{}, PartitionKey, RangeKey comparable] interface {
	Save(ctx context.Context, entity Entity) error
}

type DefaultRepository[Entity Tabler, PartitionKey, RangeKey comparable] struct {
	client *dynamodb.Client
}

func NewDefaultRepository[Entity Tabler, PartitionKey, RangeKey comparable](client *dynamodb.Client) *DefaultRepository[Entity, PartitionKey, RangeKey] {
	return &DefaultRepository[Entity, PartitionKey, RangeKey]{client: client}
}

func (d *DefaultRepository[Entity, PartitionKey, RangeKey]) Save(ctx context.Context, entity Entity) error {
	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	putItem := &dynamodb.PutItemInput{
		TableName: aws.String(entity.TableName()),
		Item:      item,
	}

	_, err = d.client.PutItem(ctx, putItem)
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
