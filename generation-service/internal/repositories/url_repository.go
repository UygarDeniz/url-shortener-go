package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	appErrors "github.com/uygardeniz/url-shortening/generation-service/internal/errors"
	"github.com/uygardeniz/url-shortening/generation-service/internal/types"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(tableName string) (*DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))

	if err != nil {
		log.Panicf("Fail to load AWS SDK config: %v", err)
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}, nil

}

func (r *DynamoDBRepository) StoreURL(ctx context.Context, url *types.URL) error {
	item, err := attributevalue.MarshalMap(url)

	if err != nil {
		log.Printf("Failed to marshal URL: %v", err)
		return fmt.Errorf("failed to marshal URL: %v", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ShortCode)"),
	})

	if err != nil {
		log.Printf("Failed to store URL: %v", err)
		var conditionalErr *dynamodbtypes.ConditionalCheckFailedException

		if errors.As(err, &conditionalErr) {
			log.Printf("Short code already exists: %s", url.ShortCode)
			return appErrors.ErrDuplicateShortCode
		}

		return fmt.Errorf("%w: %v", appErrors.ErrDatabaseAccess, err)
	}

	return nil
}
