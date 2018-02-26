package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session) // Create DynamoDB client
	}
}

func UpdateTodo(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("UpdateTodo")

	// Parse id from request body
	var (
		id        = request.PathParameters["id"]
		tableName = aws.String(os.Getenv("TODOS_TABLE_NAME"))
		done      = "done"
	)

	// Update row
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression: aws.String("set #d = :d"),
		ExpressionAttributeNames: map[string]*string{
			"#d": &done,
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {
				BOOL: aws.Bool(true),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		TableName:    tableName,
	}
	_, err := ddb.UpdateItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		return events.APIGatewayProxyResponse{ // Success HTTP response
			Body:       request.Body,
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(UpdateTodo)
}
