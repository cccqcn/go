package mysqldb

import (
	"database/sql"
	"fmt"
	"retention"
	"strconv"
	"strings"
)
import _ "github.com/go-sql-driver/mysql"

func GetRows(dsn string, sqlstr string) (*sql.Rows, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	err = db.Ping()
	retention.Trace(err)
	retention.Trace(sqlstr)
	rows, err := db.Query(sqlstr)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return rows, err
}
func InsertRows(dsn string, sqlstr string, results []*retention.Result) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	err = db.Ping()
	retention.Trace(err)
	for _, r := range results {
		ddate := r.Date.Format("2006-01-02 15:04:05")
		ddate = strings.Split(ddate, " ")[0]
		sqlone := strings.Replace(sqlstr, "@0", "'"+ddate+"'", 1)
		for i := 1; i <= 10; i++ {
			sqlone = strings.Replace(sqlone, "@"+strconv.Itoa(i), strconv.Itoa(r.Daycnts[i-1].Cnt), 1)
		}
		fmt.Println(sqlone)
		_, err := db.Exec(sqlone)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}
