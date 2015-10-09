package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var clients []Client
var random int
var ranrange int
var startServerTime string

/**
启动一个http服务器，用于生成一个随机数，并给每个访问的电脑按ip分配一个随机数，并和初始随机数计算差值并排序。
启动时候可以传参数（端口，随机数范围），默认参数（9001，100）
*/
func main() {
	flag.Parse()
	startServerTime = time.Now().Format("2006-01-02 15:04:05")
	port := "9001"
	ranrange = 100
	if flag.NArg() == 1 {
		port = flag.Arg(0)
	}
	if flag.NArg() == 2 {
		ranrange, _ = strconv.Atoi(flag.Arg(1))
	}
	//clients = make([]map[string]interface{}, 0)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	random = r1.Intn(ranrange)
	fmt.Println("Start Listening Port: " + port)
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

type Client struct {
	IP     string
	value  int
	ipTime string
}
type ClientSlice []Client

func (a ClientSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a ClientSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a ClientSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return math.Abs(float64(a[j].value-random)) > math.Abs(float64(a[i].value-random))
}

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>WHO</title>
<body><h3></h3>
`
	tableTop = `<table border="1">
<tr><th colspan="1">IP</th><th colspan="1">随机数</th><th>差值</th><th>时间</th></tr>
`

	tableItem = `
<tr><td>%v</td><td align=center>%d</td><td align=center>%d</td><td align=center>%v</td></tr>
`

	tableBottom = `
</table>`

	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

func homePage(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() // Must be called before writing response

	fmt.Fprint(writer, pageTop)
	fmt.Fprint(writer, "<p>服务器启动时间：", startServerTime, "</p>")
	fmt.Fprint(writer, "<p>系统随机数：", random, "</p>")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var myrandom int
	ip, _, _ := net.SplitHostPort(request.RemoteAddr)
	myclient, has := getIp(ip)
	if has == false {
		myrandom = r1.Intn(ranrange)
		var msg2 Client
		msg2.IP = ip
		msg2.value = myrandom
		msg2.ipTime = time.Now().Format("2006-01-02 15:04:05")
		clients = append(clients, msg2)
		sort.Sort(ClientSlice(clients))
	} else {
		myrandom = myclient.value
	}
	fmt.Fprint(writer, "<br />")
	fmt.Fprintf(writer, "本机IP："+ip+"\n\n")
	fmt.Fprint(writer, "<br />")
	fmt.Fprintf(writer, "本机随机数："+strconv.Itoa(myrandom)+"\n\n")
	fmt.Fprint(writer, "<br />")
	//fmt.Fprintf(writer, "X-Forwarded-For :"+request.Header.Get("X-FORWARDED-FOR"))
	fmt.Fprint(writer, "<br />")
	fmt.Fprint(writer, tableTop)
	for _, value := range clients {
		fmt.Fprint(writer, formatStats(value))
	}
	fmt.Fprint(writer, tableBottom)
	fmt.Fprint(writer, "<br />")

	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
	var numbers []float64
	if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
		text := strings.Replace(slice[0], ",", " ", -1)
		for _, field := range strings.Fields(text) {
			if x, err := strconv.ParseFloat(field, 64); err != nil {
				return numbers, "'" + field + "' is invalid", false
			} else {
				numbers = append(numbers, x)
			}
		}
	}
	if len(numbers) == 0 {
		return numbers, "", false // no data first time form is shown
	}
	return numbers, "", true
}

func formatStats(c Client) string {
	return fmt.Sprintf(tableItem, c.IP, c.value, random-c.value, c.ipTime)
}

func getIp(ip string) (Client, bool) {
	var c Client
	for _, value := range clients {
		if value.IP == ip {
			return value, true
		}
	}
	return c, false
}
