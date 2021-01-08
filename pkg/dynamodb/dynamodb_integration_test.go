package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = Config{
	Region:   "ap-southeast-1",
	Endpoint: "http://dynamodb:8000",
}

func TestIntegrationDynamodbGetItem_PassExisingItemKey_ReturnsItem(t *testing.T) {
	// Assemble
	key := "a"
	table := InitService(config)

	expectedItem := Item{
		Key:   key,
		Value: "ahsfkahkfahfsla",
	}

	// Act
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, actualItem, "Expected non-nil item")
	assert.Equal(t, expectedItem.Key, (*actualItem).Key)
	assert.Equal(t, expectedItem.Value, (*actualItem).Value)
}

func TestIntegrationDynamodbGetItem_PassNonExisingItemKey_ReturnsNil(t *testing.T) {
	// Assemble
	key := "x"
	table := InitService(config)

	// Act
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.Nil(t, actualItem)
}

func TestIntegrationDynamodbPutItem_PassLegalItemKey_PutsItem(t *testing.T) {
	// Assemble
	key := "w"
	table := InitService(config)

	expectedItem := Item{
		Key:   key,
		Value: "yada yada",
	}

	// Act
	table.PutItem(expectedItem.Key, expectedItem.Value)
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, actualItem, "Expected non-nil item")
	assert.Equal(t, expectedItem.Key, (*actualItem).Key)
	assert.Equal(t, expectedItem.Value, (*actualItem).Value)
}

func TestIntegrationDynamodbDeleteItem_PassExistingItemKey_DeletesItem(t *testing.T) {
	// Assemble
	key := "w"
	table := InitService(config)

	expectedItem := Item{
		Key:   key,
		Value: "yada yada",
	}

	// Act
	table.PutItem(expectedItem.Key, expectedItem.Value)
	table.DeleteItem(key)
	actualItem, err := table.GetItem(key)

	// Assert
	assert.Nil(t, err, "Expected no error")
	assert.Nil(t, actualItem)
}
