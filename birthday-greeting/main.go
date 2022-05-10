package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"birthday-greeting/dao"
	"birthday-greeting/types"
	"birthday-greeting/utils"
)

var (
	dynamoDBClient dynamodbiface.DynamoDBAPI

	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"
	ErrNoIP               = errors.New("No IP in HTTP response")
	ErrNon200Response     = errors.New("Non 200 Response found")
	TableName             = "user"
	GlobalSecondaryIndex  = "birthMonth-birthDay-index"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// validate the APIGateway Request
	apigatewayResp, err := validateAPIGatewayRequest()
	if err != nil {
		fmt.Printf("invalid APIGatewayProxyRequest, %v\n", err)
		return apigatewayResp, err
	}

	// Make the query to Database
	res, err := dao.QueryByGSI(dynamoDBClient, TableName, GlobalSecondaryIndex)
	if err != nil {
		fmt.Printf("failed to QueryByGSI, %v\n", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// Marshal to JSON object
	// e.g. dynamoDB item -> user -> create greeting message -> Marshal to json format
	var greetingList []types.BirthdayGreeting
	for _, item := range res.Items {
		user := types.User{}
		dynamodbattribute.UnmarshalMap(item, &user)
		greeting, err := utils.CraftBirthdayGreetingForUser(user)
		if err != nil {
			fmt.Printf("Failed to craft birthday greeting, ignore this user: %v\n", err)
			continue
		}
		greetingList = append(greetingList, greeting)
	}
	greetingsInJsonFmt, _ := json.Marshal(greetingList)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(string(greetingsInJsonFmt)), // put json here
		StatusCode: 200,
	}, nil
}

func validateAPIGatewayRequest() (events.APIGatewayProxyResponse, error) {
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
	return events.APIGatewayProxyResponse{}, err
}

func main() {
	dynamoDBClient = dao.NewDynamoDBClient() // singleton
	lambda.Start(handler)
}
