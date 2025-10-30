package oldfart

import (
	"context"
	"oldfartscounter/internal/aws/dynamodb"

	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBRepository struct {
	repo dynamodb.Repository[*OldFart, SteamId, dynamodb.NilPointer]
}

func (d DynamoDBRepository) Save(ctx context.Context, entity *OldFart) error {
	return d.repo.Save(ctx, entity)
}

func NewDynamoDBRepository(client *awsdynamodb.Client) *DynamoDBRepository {
	return &DynamoDBRepository{
		repo: dynamodb.NewDefaultRepository[*OldFart, SteamId, dynamodb.NilPointer](client),
	}
}
