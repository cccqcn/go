package main

import (
	"database/sql"
	"fmt"
	//	"log"
	"strconv"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type User struct {
	Pid   int
	Dates []time.Time
}

var users []*User

func main() {
	db, err := sql.Open("mysql", "")
	if err != nil {
		fmt.Println(err)
	}
	err = db.Ping()
	fmt.Println(err)

	rows, err := db.Query("SELECT * FROM ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(rows)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	cnt := 0
	users = make([]*User, 0)
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		v, _ := strconv.ParseInt(string(values[3]), 10, 64)
		t := time.Unix(v/1000, 0)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if columns[i] == "time" {
				value = t.Format(time.RFC3339)
			}
			fmt.Println(columns[i], ": ", value)
		}
		pid, _ := strconv.Atoi(string(values[2]))
		addUser(pid, t)
		fmt.Println("-----------------------------------")
		cnt++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(cnt)
	for _, user := range users {
		dd := ""
		for _, date := range user.Dates {
			dd += date.Month().String() + "." + strconv.Itoa(date.Day()) + ","
		}
		fmt.Println(user.Pid, user.Dates)
	}
	fmt.Println(len(users))
}
func addUser(pid int, t time.Time) {
	olduser := getUserByPid(pid)
	if olduser.Pid == 0 {
		ts := make([]time.Time, 0)
		ts = append(ts, t)
		user := User{Pid: pid, Dates: ts}
		users = append(users, &user)
	} else {
		for _, d := range olduser.Dates {
			if d.Year() == t.Year() && d.Month() == t.Month() && d.Day() == t.Day() {
				return
			}
		}
		//		(&olduser).AddTime(t)
		olduser.Dates = append(olduser.Dates, t)
		fmt.Println("append", olduser.Pid, pid, t.Year(), t.Month(), t.Day(), len(olduser.Dates))
	}
}
func getUserByPid(pid int) (u *User) {
	for _, user := range users {
		if user.Pid == pid {
			return user
		}
	}
	return &User{Pid: 0}
}
