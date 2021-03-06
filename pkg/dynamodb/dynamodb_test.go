package dynamodb

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gusaul/go-dynamock"
	"github.com/stretchr/testify/assert"
)

func unmarshalItem(marshalledItem dynamodb.GetItemOutput) (item *Item, err error) {
	unmarshalledItem := Item{}
	err = dynamodbattribute.UnmarshalMap(marshalledItem.Item, &unmarshalledItem)

	return &unmarshalledItem, err
}

func TestDynamodbGetItem_PassExisingItemKey_ReturnsItem(t *testing.T) {
	mockedDb, mock := dynamock.New()

	// Assemble
	key := "a"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}
	expectedKey := map[string]*dynamodb.AttributeValue{
		"Key": {S: aws.String(key)},
	}
	expectedMarshalledItem := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"Key":   {S: aws.String(key)},
			"Value": {S: aws.String("my expected value")},
		},
	}

	mock.ExpectGetItem().
		ToTable(table.tableName).
		WithKeys(expectedKey).
		WillReturns(expectedMarshalledItem)

	expectedItem, _ := unmarshalItem(expectedMarshalledItem)

	// Act
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, actualItem, "Expected non-nil item")
	assert.Equal(t, expectedItem, actualItem)
}
func TestDynamodbGetItem_PassNonExisingItemKey_ReturnsNil(t *testing.T) {
	mockedDb, mock := dynamock.New()

	// Assemble
	key := "a"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}
	expectedKey := map[string]*dynamodb.AttributeValue{
		"Key": {S: aws.String(key)},
	}
	expectedMarshalledItem := dynamodb.GetItemOutput{
		Item: nil,
	}

	mock.ExpectGetItem().
		ToTable(table.tableName).
		WithKeys(expectedKey).
		WillReturns(expectedMarshalledItem)

	// Act
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.Nil(t, actualItem)
}

func TestDynamodbGetItem_PassNonExisingItemKey_GetError(t *testing.T) {
	mockedDb, _ := dynamock.New()

	// Assemble
	key := "a"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}

	// Act
	_, err := table.GetItem(key)

	// Assert
	assert.NotNil(t, err)
}

func TestDynamodbPutItem_PassItemKey_GetNoError(t *testing.T) {
	mockedDb, mock := dynamock.New()

	// Assemble
	key := "a"
	value := "sdfghjkjhgf"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}
	item := map[string]*dynamodb.AttributeValue{
		"Key":   {S: aws.String(key)},
		"Value": {S: aws.String(value)},
	}

	mock.ExpectPutItem().
		ToTable(table.tableName).
		WithItems(item).
		WillReturns(dynamodb.PutItemOutput{})

	// Act
	err := table.PutItem(key, value)

	// Assert
	assert.Nil(t, err)
}

func TestDynamodbPutItem_PassItemKey_GetError(t *testing.T) {
	mockedDb, _ := dynamock.New()

	// Assemble
	key := "a"
	value := "sdfghjkjhgf"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}

	// Act
	err := table.PutItem(key, value)

	// Assert
	assert.NotNil(t, err)
}

func TestDynamodbDeleteItem_PassItemKey_GetNoError(t *testing.T) {
	mockedDb, mock := dynamock.New()

	// Assemble
	key := "a"

	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}
	expectedKey := map[string]*dynamodb.AttributeValue{
		"Key": {S: aws.String(key)},
	}

	mock.ExpectDeleteItem().
		ToTable(table.tableName).
		WithKeys(expectedKey).
		WillReturns(dynamodb.DeleteItemOutput{})

	// Act
	err := table.DeleteItem(key)

	// Assert
	assert.Nil(t, err)
}

func TestDynamodbDeleteItem_PassItemKey_GetError(t *testing.T) {
	mockedDb, _ := dynamock.New()

	// Assemble
	key := "a"
	table := serviceImpl{
		db:        mockedDb,
		tableName: defaultTableName,
	}

	// Act
	err := table.DeleteItem(key)

	// Assert
	assert.NotNil(t, err)
}
