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

// Config provides abstraction over DynamoDB servie configuration
type Config struct {
	Region    string
	TableName string
	Endpoint  string
}

// Item provides abstraction over DynamoDB entry
type Item struct {
	Key   string
	Value string
}

// Service interface provides abstraction over DynamoDB table
type Service interface {
	GetItem(key string) (*Item, error)
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

// InitService ...
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
