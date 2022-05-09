package main

import (
	"encoding/json"
	// "errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"birthday-greeting/dao"
	"birthday-greeting/types"
)

var (
	dynamoDBClient dynamodbiface.DynamoDBAPI

	TableName            = "user"
	GlobalSecondaryIndex = "birthMonth-birthDay-index"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// resp, err := http.Get(DefaultHTTPGetAddress)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// if resp.StatusCode != 200 {
	// 	return events.APIGatewayProxyResponse{}, ErrNon200Response
	// }

	// ip, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// if len(ip) == 0 {
	// 	return events.APIGatewayProxyResponse{}, ErrNoIP
	// }

	res, err := dao.QueryByGSI(dynamoDBClient, TableName, GlobalSecondaryIndex)
	if err != nil {
		fmt.Printf("failed to QueryByGSI, %v\n", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// dynamoDB item -> user class -> Marshal to json format
	var userList []types.User
	for _, item := range res.Items {
		user := types.User{}
		dynamodbattribute.UnmarshalMap(item, &user)
		userList = append(userList, user)
	}
	usersInJsonFmt, _ := json.Marshal(userList)

	return events.APIGatewayProxyResponse{
		// put json here
		Body:       fmt.Sprintf(string(usersInJsonFmt)),
		StatusCode: 200,
	}, nil
}

func main() {
	dynamoDBClient = dao.NewDynamoDBClient()
	lambda.Start(handler)
}
