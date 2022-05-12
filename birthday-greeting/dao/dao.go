package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"birthday-greeting/secrets"
	"birthday-greeting/types"
	"birthday-greeting/utils"
)

var (
	rdsSecret secrets.Secret
)

func init() {
	rdsSecret, _ = secrets.GetSecret()
}

func GetRDSDB(dbName string) (*sql.DB, error) {
	driver, user, password, endpoint, port, dbName := rdsSecret.Engine, rdsSecret.UserName, rdsSecret.Password, rdsSecret.Host, rdsSecret.Port, rdsSecret.DBName
	charset := "charset=utf8mb4"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?%s", user, password, endpoint, port, dbName, charset)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		fmt.Printf("unable to connect to DB: %v\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("unable to ping the DB: %v\n", err)
		return nil, err
	}
	return db, nil
}

func GetUsersToGreet(db *sql.DB, tableName string) ([]types.BirthdayGreeting, error) {

	mmddyyyy := time.Now().Format("01-02-2006")
	date, _ := strconv.Atoi(mmddyyyy[3:5])
	month, _ := strconv.Atoi(mmddyyyy[:2])
	queryString := "SELECT * FROM " + tableName + " WHERE birthday=" + strconv.Itoa(date) + " AND birthmonth=" + strconv.Itoa(month) + ";"

	fmt.Printf("[debug] queryString: %v\n", queryString)
	rows, err := db.Query(queryString)
	if err != nil {
		fmt.Printf("error Query on table %v, :%v", tableName, err)
	}

	// e.g. RDS item -> user -> create greeting message -> Marshal to json format
	var greetingList []types.BirthdayGreeting
	for rows.Next() {
		var id, firstname, lastname, email, gender, dateofbirth string
		var dd, mm, yyyy int

		err := rows.Scan(&id, &gender, &email, &firstname, &lastname, &dateofbirth, &dd, &mm, &yyyy)
		if err != nil {
			fmt.Printf("error scan: %v", err)
		}
		fmt.Println("user: ", id, firstname, lastname, email, gender, dateofbirth)
		user := types.User{
			FirstName: firstname,
		}
		greeting, err := utils.CraftBirthdayGreetingForUser(user)
		if err != nil {
			fmt.Printf("Failed to craft birthday greeting, ignore this user: %v\n", err)
			continue
		}
		greetingList = append(greetingList, greeting)
	}
	return greetingList, nil
}
