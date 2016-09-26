package retention

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var users []*User
var dates []time.Time
var results []*Result

func Results() []*Result {
	return results
}
func Output() {
	Trace("totalUsers: ", len(users))
	Trace("totalDates: ", len(dates))
	Trace("totalResults: ", len(results))
	for _, rr := range results {
		dcstr := ""
		for _, dc := range rr.Daycnts {
			dcstr += strconv.Itoa(dc.Days) + ":" + strconv.Itoa(dc.Cnt) + ","
		}
		datetime := rr.Date.Format("2006-01-02 15:04:05")
		date := strings.Split(datetime, " ")[0]
		fmt.Println(date, dcstr)
	}
}
func AnalysisRows(rows *sql.Rows) {
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
	dates = make([]time.Time, 0)

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
		addDate(&dates, t)
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
			Trace(columns[i], ": ", value)
		}
		pid, _ := strconv.Atoi(string(values[2]))
		addUser(pid, t)
		Trace("-----------------------------------")
		cnt++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	Trace(cnt)
	for _, user := range users {
		dd := ""
		for _, date := range user.Dates {
			dd += date.Month().String() + "." + strconv.Itoa(date.Day()) + ","
		}
		Trace(user.Pid, dd)
	}

	results = make([]*Result, 0)
	for _, d := range dates {
		newr := newResult(d)
		results = append(results, &newr)
		for _, u := range users {
			if isRegisterOnDay(*u, d) {
				for _, dc := range newr.Daycnts {
					dAfterDays := d.AddDate(0, 0, dc.Days)
					isLogin := isLoginOnDays(*u, dAfterDays)
					if isLogin {
						dc.Cnt++
					}
				}
			}
		}
	}
}

func newResult(d time.Time) Result {
	cnts := make([]*DayCnt, 0)
	cnts = append(cnts, &DayCnt{Days: 0, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 1, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 2, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 3, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 4, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 5, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 6, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 7, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 14, Cnt: 0})
	cnts = append(cnts, &DayCnt{Days: 30, Cnt: 0})
	newr := Result{Date: d, Daycnts: cnts}
	return newr
}
func isRegisterOnDay(u User, t time.Time) bool {
	d := u.Dates[0]
	if d.Year() == t.Year() && d.Month() == t.Month() && d.Day() == t.Day() {
		return true
	}
	return false
}
func isLoginOnDays(u User, t time.Time) bool {
	for _, d := range u.Dates {
		if d.Year() == t.Year() && d.Month() == t.Month() && d.Day() == t.Day() {
			return true
		}
	}
	return false
}
func addUser(pid int, t time.Time) {
	olduser := getUserByPid(pid)
	if olduser.Pid == 0 {
		ts := make([]time.Time, 0)
		addDate(&ts, t)
		user := User{Pid: pid, Dates: ts}
		users = append(users, &user)
	} else {
		addDate(&(olduser.Dates), t)
		Trace("append", olduser.Pid, pid, t.Year(), t.Month(), t.Day(), len(olduser.Dates))
	}
}
func addDate(ds *([]time.Time), t time.Time) {
	for _, d := range *ds {
		if d.Year() == t.Year() && d.Month() == t.Month() && d.Day() == t.Day() {
			return
		}
	}
	tt, _ := time.Parse(timeFormat, t.Format(timeFormat))
	Trace(tt)
	*ds = append(*ds, tt)
}
func getUserByPid(pid int) (u *User) {
	for _, user := range users {
		if user.Pid == pid {
			return user
		}
	}
	return &User{Pid: 0}
}
func Trace(a ...interface{}) {
	if TraceFlag == true {
		fmt.Println(a)
	}
}
