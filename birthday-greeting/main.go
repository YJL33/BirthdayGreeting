package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"birthday-greeting/dao"
)

var (
	dynamoDBClient dynamodbiface.DynamoDBAPI

	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"
	ErrNoIP               = errors.New("No IP in HTTP response")
	ErrNon200Response     = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	res, err := dao.QueryByGSI(dynamoDBClient, "user", "birthMonth-birthDay-index")
	if err != nil {
		fmt.Printf("failed to QueryByGSI, %v\n", err)
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Sprintf("log?, %v", res.String())

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v, %v", string(ip), res.String()),
		StatusCode: 200,
	}, nil
}

func main() {
	dynamoDBClient = dao.NewDynamoDBClient()
	lambda.Start(handler)
}
