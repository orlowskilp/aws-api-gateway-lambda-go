package dynamodb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	defaultRegion    string = "ap-southeast-1"
	defaultTableName string = "sample-table"
	keyNameString    string = "Key"
)

type serviceImpl struct {
	tableName string
	db        dynamodbiface.DynamoDBAPI
}

// Config provides abstraction over DynamoDB service configuration.
// It controls the following parameters:
// * Region    - AWS region. If not specified in ~/.aws/config file
//               the region is set to `defaultRegion` constant
// * TableName - Name of DynamoDB table. If not set explicitly, the
//               service will be initialized with value set in the
//               `defaultTableName` constant
// * Endpoint  - Only applicable for testing using a local instance
//               of DynamoDB e.g. https://dynamodb:8000. Leave blank
//               to operate on a table in AWS.
type Config struct {
	Region    string
	TableName string
	Endpoint  string
}

// Item provides abstraction over DynamoDB item (entry).
// * Key   - Represents index of the item
// * Value - Represents value assigned to the index in `Key`
type Item struct {
	Key   string
	Value string
}

// Service interface provides CRUD abstraction over DynamoDB table
type Service interface {
	GetItem(key string) (*Item, error)
	PutItem(key string, value string) error
	DeleteItem(key string) error
}

func initContext(config Config) *dynamodb.DynamoDB {
	awsConfig := aws.Config{}

	if len(config.Region) > 0 {
		awsConfig.Region = aws.String(config.Region)
	} else {
		awsConfig.Region = aws.String(defaultRegion)
	}

	if len(config.Region) > 0 {
		awsConfig.Endpoint = aws.String(config.Endpoint)
	}

	sesh := session.Must(session.NewSession(&awsConfig))

	return dynamodb.New(sesh)
}

// InitService initializes the service with parameters supplied
// in the `config` parameter
func InitService(config Config) Service {
	tabName := config.TableName
	if len(tabName) <= 0 {
		tabName = defaultTableName
	}

	context := serviceImpl{
		tableName: tabName,
		db:        initContext(config),
	}

	return context
}

func (service serviceImpl) GetItem(key string) (*Item, error) {
	result, err := service.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyNameString: {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		log.Printf("Couldn't find item with key: %s", key)
		return nil, nil
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &item, err
}

func (service serviceImpl) PutItem(key string, value string) error {
	item := Item{
		Key:   key,
		Value: value,
	}

	marshalledItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal Record, %v", err))
	}

	_, err = service.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(service.tableName),
		Item:      marshalledItem,
	})

	if err != nil {
		return err
	}

	log.Printf("Successfully put item with key: %s", key)

	return nil
}

func (service serviceImpl) DeleteItem(key string) error {
	_, err := service.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyNameString: {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return err
	}

	log.Printf("Successfully deleted item with key: %s", key)

	return nil
}
