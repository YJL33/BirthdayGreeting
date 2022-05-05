package dao

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoDBInstance  dynamodbiface.DynamoDBAPI
	dynamoDBSingleton sync.Once
	region            string
	sess              *session.Session
)

func init() {
	region = os.Getenv("AWS_REGION")
	sess = session.Must(session.NewSession())
}

// NewDynamoDBClient returns a new dynamoDB client
func NewDynamoDBClient() dynamodbiface.DynamoDBAPI {
	dynamoDBSingleton.Do(func() {
		dynamoDBInstance = dynamodb.New(sess, &aws.Config{Region: &region})
	})
	return dynamoDBInstance
}

// // GetItemByPrimaryKey simply gets item by primary key
// func GetItemByPrimaryKey(dao dynamodbiface.DynamoDBAPI, tableName string, primaryKey string, primaryKeyValue string) (*dynamodb.GetItemOutput, error) {
// 	fmt.Printf("getItem from table: %s, key: %s\n", tableName, primaryKeyValue)
// 	provision, err := dao.GetItem(&dynamodb.GetItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			primaryKey: {
// 				S: aws.String(primaryKeyValue),
// 			},
// 		},
// 		TableName: aws.String(tableName),
// 	})
// 	if err != nil {
// 		fmt.Printf("failed to make getItems API call, %v\n", err)
// 		return nil, err
// 	}

// 	return provision, nil
// }
