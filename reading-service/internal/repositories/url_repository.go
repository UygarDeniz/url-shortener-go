package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/uygardeniz/url-shortening/reading-service/internal/errors"
	"github.com/uygardeniz/url-shortening/reading-service/internal/types"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(tableName string) (*DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}, nil

}

func (r *DynamoDBRepository) GetOriginalURL(ctx context.Context, shortCode string) (*types.URL, error) {

	key := map[string]dynamodbtypes.AttributeValue{
		"ShortCode": &dynamodbtypes.AttributeValueMemberS{Value: shortCode},
	}

	res, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key:       key,
	})

	if err != nil {
		log.Printf("Error getting item from DynamoDB: %v", err)
		return nil, errors.ErrDatabaseAccess
	}

	if res.Item == nil {
		return nil, errors.ErrURLNotFound
	}

	var url types.URL
	err = attributevalue.UnmarshalMap(res.Item, &url)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDataFormat, err)
	}

	return &url, nil

}
