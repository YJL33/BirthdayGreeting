package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"birthday-greeting/dao"
)

var (
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"
	ErrNoIP               = errors.New("No IP in HTTP response")
	ErrNon200Response     = errors.New("Non 200 Response found")
	DBName                = "userDB"
	TableName             = "user"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// validate the APIGateway Request
	apigatewayResp, err := validateAPIGatewayRequest()
	if err != nil {
		fmt.Printf("invalid APIGatewayProxyRequest, %v\n", err)
		return apigatewayResp, err
	}

	// connect to DB
	rdsDB, err := dao.GetRDSDB(DBName)
	if err != nil {
		fmt.Printf("failed to get RDS DB, %v\n", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// Get the users to greet
	greetingList, err := dao.GetUsersToGreet(rdsDB, TableName)
	if err != nil {
		fmt.Printf("failed to get Users to Greet, %v\n", err)
		return events.APIGatewayProxyResponse{}, err
	}
	greetingsInJsonFmt, _ := json.Marshal(greetingList)
	fmt.Printf("jsonFmt: %v\n", string(greetingsInJsonFmt))

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
	lambda.Start(handler)
}
