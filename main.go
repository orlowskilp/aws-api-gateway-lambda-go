package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/orlowskilp/aws-api-gateway-lambda-go/pkg/dynamodb"
)

const (
	defaultAWSRegion         string = "ap-southeast-1"
	keyPathParameterName     string = "key"
	httpSuccessCode          int    = 200
	httpServerErrorCode      int    = 500
	tableNameEnvVariableName string = "TABLE_NAME"
)

func handleGetRequest(key string) events.APIGatewayProxyResponse {
	if len(key) <= 0 {
		message := "You need to specify /key/{key} path parameter."
		return events.APIGatewayProxyResponse{
			Body:       message,
			StatusCode: httpSuccessCode,
		}
	}

	table := dynamodb.InitService(dynamodb.Config{
		Region:    defaultAWSRegion,
		TableName: os.Getenv(tableNameEnvVariableName),
	})
	entry, err := table.GetItem(key)

	if err != nil {
		log.Fatalln(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: httpServerErrorCode,
		}
	}

	if entry == nil {
		message := "No entry found for key: %s"
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(message, key),
			StatusCode: httpSuccessCode,
		}
	}

	message := "Value for key: %s is %s"
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(message, key, (*entry).Value),
		StatusCode: httpSuccessCode,
	}
}

func handlePutRequest(key string, value string) events.APIGatewayProxyResponse {
	if len(key) <= 0 || len(value) <= 0 {
		message := "Key and value parameters must not be empty"
		return events.APIGatewayProxyResponse{
			Body:       message,
			StatusCode: httpServerErrorCode,
		}
	}

	table := dynamodb.InitService(dynamodb.Config{
		Region: defaultAWSRegion,
	})

	err := table.PutItem(key, value)

	if err != nil {
		log.Fatalln(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: httpServerErrorCode,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: httpSuccessCode,
	}
}

func handleDeleteRequest(key string) events.APIGatewayProxyResponse {
	if len(key) <= 0 {
		message := "Key parameter must not be empty"
		return events.APIGatewayProxyResponse{
			Body:       message,
			StatusCode: httpServerErrorCode,
		}
	}

	table := dynamodb.InitService(dynamodb.Config{
		Region: defaultAWSRegion,
	})

	err := table.DeleteItem(key)

	if err != nil {
		log.Fatalln(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: httpServerErrorCode,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: httpSuccessCode,
	}
}

func handleRequest(ctx context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := request.PathParameters[keyPathParameterName]
	value := request.Body
	httpMethod := request.HTTPMethod

	var response events.APIGatewayProxyResponse

	switch httpMethod {
	case "GET":
		response = handleGetRequest(key)
		break
	case "PUT":
		response = handlePutRequest(key, value)
		break
	case "DELETE":
		response = handleDeleteRequest(key)
		break
	}

	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
