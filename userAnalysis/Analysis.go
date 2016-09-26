package main

import (
	"fmt"
	//	"log"
	"flag"
	"mysqldb"
	"retention"
	//	"strings"

	"github.com/Unknwon/goconfig"
)

var cfg *goconfig.ConfigFile

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("no args")
		return
	}
	config := flag.Arg(0)

	var err error
	cfg, err = goconfig.LoadConfigFile(config)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	retention.TraceFlag = cfg.MustBool(goconfig.DEFAULT_SECTION, "traceFlag")

	dsn, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "dsn")
	sqlstr, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "sql")
	retention.Trace(sqlstr)
	rows, err := mysqldb.GetRows(dsn, sqlstr)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	retention.Trace(rows)
	retention.AnalysisRows(rows)
	retention.Output()
	dsn2, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "dsn2")
	sqlstr2, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "sql2")
	mysqldb.InsertRows(dsn2, sqlstr2, retention.Results())
}
