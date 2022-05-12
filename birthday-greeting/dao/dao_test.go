package dao

import (
	"birthday-greeting/types"
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetRDSDB(t *testing.T) {
	type args struct {
		dbName string
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRDSDB(tt.args.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRDSDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRDSDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsersToGreet(t *testing.T) {
	type args struct {
		db        *sql.DB
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []types.BirthdayGreeting
		wantErr bool
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUsersToGreet(tt.args.db, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsersToGreet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsersToGreet() = %v, want %v", got, tt.want)
			}
		})
	}
}
