package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	ErrBodyNotProvided   = errors.New("no HTTP body")
	ErrUnmarshallProblem = errors.New("error unmarshalling payload")
	ErrMandatoryElementsMissing = errors.New("Input must provide both instance and notify fields")
)

type SubRequest struct {
	InstanceID string `json:"instance"`
	Notify     string `json:"notify"`
}

type SubResponse struct {
	Confirmation string
}

type AWSContext struct {
	ddbSvc dynamodbiface.DynamoDBAPI
}

var subscriptionTable = os.Getenv("SUBSCRIPTION_TABLE")

func writeSubscriptionInfo(ddbSvc dynamodbiface.DynamoDBAPI, subscribe *SubRequest) error {

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"InstanceID": {
				S: aws.String(subscribe.InstanceID),
			},
			"Notify": {
				S: aws.String(subscribe.Notify),
			},
		},
		TableName: aws.String(subscriptionTable),
	}
	_, err := ddbSvc.PutItem(input)

	return err
}

func makeHandler(awsContext *AWSContext) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Handler is your Lambda function handler
	// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
	// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

		log.Printf("Request body: %s\n", request.Body)

		//Unmarshall request
		var subscribe SubRequest
		if err := json.Unmarshal([]byte(request.Body), &subscribe); err != nil {
			log.Println("error unmarshalling request")
			return events.APIGatewayProxyResponse{}, ErrUnmarshallProblem
		}

		if subscribe.InstanceID == "" || subscribe.Notify == "" {
			log.Println("inputs not fully specified")
			return events.APIGatewayProxyResponse{}, ErrMandatoryElementsMissing
		}

		if err := writeSubscriptionInfo(awsContext.ddbSvc, &subscribe); err != nil {
			log.Println("error persisting subscription information")
			return events.APIGatewayProxyResponse{}, err
		}

		// If no name is provided in the HTTP request body, throw an error
		if len(request.Body) < 1 {
			return events.APIGatewayProxyResponse{}, ErrBodyNotProvided
		}

		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 200,
		}, nil

	}
}

func main() {
	var awsContext AWSContext

	sess := session.New()
	svc := dynamodb.New(sess)

	awsContext.ddbSvc = svc

	handler := makeHandler(&awsContext)
	lambda.Start(handler)
}
