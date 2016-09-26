package retention

import (
	"time"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02"
)

var TraceFlag bool

type User struct {
	Pid   int
	Dates []time.Time
}
type DayCnt struct {
	Days int
	Cnt  int
}
type Result struct {
	Date    time.Time
	Daycnts []*DayCnt
}
