package dao

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	mockOutput       *dynamodb.QueryOutput
	defaultTableName string
	defaultGSI       string
)

func init() {
	mockOutput = &dynamodb.QueryOutput{}
	defaultTableName = "user"
	defaultGSI = "birthMonth-birthDay-index"
}

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if *input.TableName != defaultTableName || *input.IndexName != defaultGSI {
		fmt.Printf("TableName:, %v\n", *input.TableName)
		fmt.Printf("defaultTableName:, %v\n", defaultTableName)
		fmt.Printf("IndexName:, %v\n", *input.IndexName)
		fmt.Printf("defaultGSI:, %v\n", defaultGSI)
		return nil, errors.New("invalid tableName or GSI")
	}
	return mockOutput, nil
}

func TestQueryByGSI(t *testing.T) {
	// Setup Test
	mockSvc := &mockDynamoDBClient{}

	type args struct {
		dao       dynamodbiface.DynamoDBAPI
		tableName string
		gsi       string
	}

	tests := []struct {
		name    string
		args    args
		want    *dynamodb.QueryOutput
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				dao:       mockSvc,
				tableName: defaultTableName,
				gsi:       defaultGSI,
			},
			want:    mockOutput,
			wantErr: false,
		},
		{
			name: "Invalid table name",
			args: args{
				dao:       mockSvc,
				tableName: "xxxx",
				gsi:       defaultGSI,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid GSI name",
			args: args{
				dao:       mockSvc,
				tableName: defaultTableName,
				gsi:       "xxxx",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryByGSI(tt.args.dao, tt.args.tableName, tt.args.gsi)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryByGSI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByGSI() = %v, want %v", got, tt.want)
			}
		})
	}
}
