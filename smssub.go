package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

var (
	ErrBodyNotProvided   = errors.New("no HTTP body")
	ErrUnmarshallProblem = errors.New("error unmarshalling payload")
)

type SubRequest struct {
	InstanceID string `json:instance`
	Notify     string `json:notify`
}

type SubResponse struct {
	Confirmation string
}

type AWSContext struct {
}

var subscriptionTable = os.Getenv("SUBSCRIPTION_TABLE")

func writeSubscriptionInfo(subscribe *SubRequest) error {
	return nil
}

func makeHandler(awsContent *AWSContext) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Handler is your Lambda function handler
	// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
	// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

		//Unmarshall request
		var subscribe SubRequest
		if err := json.Unmarshal([]byte(request.Body), &subscribe); err != nil {
			log.Println("error unmarshalling request")
			return events.APIGatewayProxyResponse{}, ErrUnmarshallProblem
		}

		if err := writeSubscriptionInfo(&subscribe); err != nil {
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
	var awsContent AWSContext
	handler := makeHandler(&awsContent)
	lambda.Start(handler)
}
