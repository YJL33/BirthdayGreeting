package dao

import (
	"fmt"
	"os"
	"sync"
	"time"

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

// QueryByGSI simply make a query to DynamoDB by GSI
// below is the example of command line query
// aws dynamodb query \
//     --table-name user \
//     --index-name birthMonth-birthDay-index \
//     --key-condition-expression "birthDay = :birthDayAttr AND birthMonth = :birthMonthAttr" \
//     --expression-attribute-values  '{":birthDayAttr":{"N":"8"}, ":birthMonthAttr":{"N":"8"}}'
func QueryByGSI(dao dynamodbiface.DynamoDBAPI, tableName string, gsi string) (*dynamodb.QueryOutput, error) {
	fmt.Printf("Query from table: %s, with GSI: %s\n", tableName, gsi)
	currentTime := time.Now()
	MMDDYYYY := currentTime.Format("01-02-2006")
	month := MMDDYYYY[:2]
	date := MMDDYYYY[3:5]
	results, err := dao.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":birthMonthAttr": {
				N: aws.String(month),
			},
			":birthDayAttr": {
				N: aws.String(date),
			},
		},
		IndexName:              &gsi,
		KeyConditionExpression: aws.String("birthDay = :birthDayAttr AND birthMonth = :birthMonthAttr"),
		TableName:              aws.String(tableName),
	})
	if err != nil {
		fmt.Printf("failed to make Query on AWS DynamoDB, %v\n", err)
		return nil, err
	}

	return results, nil
}
