package main

import (
	"testing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)


type dynamoDBMockery struct {
	dynamodbiface.DynamoDBAPI
}

func (m *dynamoDBMockery) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	var out dynamodb.PutItemOutput
	return &out, nil
}

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: `{"instance":"foo","notify":"+12223334444"}`, RequestContext:events.APIGatewayProxyRequestContext{RequestID:"resource1"}},
			expect:  "",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "", RequestContext:events.APIGatewayProxyRequestContext{RequestID:"resource2"}},
			expect:  "",
			err:     ErrBodyNotProvided,
		},
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: `{"what":"foo","notify":"+12223334444"}`, RequestContext:events.APIGatewayProxyRequestContext{RequestID:"resource1"}},
			expect:  "",
			err:     ErrMandatoryElementsMissing,
		},
	}

	var awsContext AWSContext
	var myMock dynamoDBMockery
	awsContext.ddbSvc = &myMock
	handler := makeHandler(&awsContext)

	for _, test := range tests {
		response, err := handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}