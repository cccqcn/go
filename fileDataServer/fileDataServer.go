package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"  
//	"math"
//	"math/rand"
//	"net"
	"net/http"
    "path"
	"os"
	"strconv"
	"strings"
	"time"
	"bufio"
	"io"
	"io/ioutil"
"golang.org/x/text/encoding/simplifiedchinese"
"code.google.com/p/mahonia"
)

var startServerTime string
var xmlpath string

/**
启动一个http服务器，用于生成一个随机数，并给每个访问的电脑按ip分配一个随机数，并和初始随机数计算差值并排序。
启动时候可以传参数（端口，随机数范围），默认参数（9001，100）
*/
func main() {
	flag.Parse()
	startServerTime = time.Now().Format("2006-01-02 15:04:05")
	port := "9001"
	if flag.NArg() == 1 {
		port = flag.Arg(0)
	}
	f, _ := os.Getwd()
	xmlpath = "E:\\Projects\\wolf\\battledata"
    fmt.Println("current path:", f)
	xmlpath = f
	if flag.NArg() == 2 {
		xmlpath = flag.Arg(1)
	}
	//clients = make([]map[string]interface{}, 0)
	fmt.Println("Start Listening Port: " + port)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/crossdomain.xml", crossdomain)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func crossdomain(writer http.ResponseWriter, request *http.Request) {
	s :=`<?xml version="1.0"?> 
	<!-- http://www.foo.com/crossdomain.xml --> 
	<cross-domain-policy> 
	<allow-access-from domain="*" /> 
	</cross-domain-policy>`
	fmt.Fprint(writer, s)
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() // Must be called before writing response

	//fmt.Fprint(writer, "<p>服务器启动时间：", startServerTime, "</p>")

	if slice, found := request.Form["cmd"]; found && len(slice) > 0 {
		cmd := slice[0]
		switch cmd {
			case "list":
				getFilelist(writer)
			case "get":
				file, _ := request.Form["file"];
				getFile(writer, file[0])
			case "save":
				file, _ := request.Form["file"];
				xml, _ := request.Form["xml"];
				saveFile(writer, file[0], xml[0])
		}
	}
	if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
		fmt.Fprint(writer, slice[0])
	}
}

func saveFile(writer http.ResponseWriter, file string, xml string) { 
	file, _ = GBKToUtf8(file) // 文件名转换为 GBK编码
	xml, _ = GBKToUtf8(xml) // 文件名转换为 GBK编码
	fullpath := xmlpath + "\\" + file + ".xml"
    fmt.Println("save", file)
	_, err := os.Stat(fullpath)
	if err != nil {
		fin, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0644)
        defer fin.Close()
        if err != nil {//文件打开错误
            fmt.Println(fullpath, err)
			return
        }
		fwriter := bufio.NewWriter(fin)
		defer func() {
			if err == nil {
				err = fwriter.Flush()
			}
		}()
		if _, err = fwriter.WriteString(xml); err != nil {//重新写入文件
	        fmt.Println(err)
			return 
		}
		fmt.Fprint(writer, "OK")
	} else{
		fmt.Fprint(writer, "文件重复")
	}
}

func getFile(writer http.ResponseWriter, file string) {  
	file, _ = GBKToUtf8(file) // 文件名转换为 GBK编码
	fullpath := xmlpath + "\\" + file + ".xml"
    fmt.Println("get", fullpath)
	fileBytes, err := ioutil.ReadFile(fullpath)
	if err != nil || io.EOF == err {//读取失败
        fmt.Println(err)
		return
	}
	line := string(fileBytes)
	fmt.Fprint(writer, line)
}  

func getFilelist(writer http.ResponseWriter) {  
    err := filepath.Walk(xmlpath, func(pathstr string, f os.FileInfo, err error) error {  
        if f == nil {  
            return err  
        }  
        if f.IsDir() {  
            return nil  
        }
		pathstr = strings.Replace(pathstr, "\\", "/", -1)
    	filenameWithSuffix := path.Base(pathstr)

		fileSuffix := path.Ext(filenameWithSuffix)
		if fileSuffix == ".xml"{
   			filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			fmt.Fprint(writer, filenameOnly + "<br />")
		}
        return nil  
        })  
    if err != nil {  
        fmt.Printf("filepath.Walk() returned %v\n", err)  
    }  
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
func convert(s string)string {
    var dec mahonia.Decoder
    dec = mahonia.NewDecoder("gbk")
    if ret, ok := dec.ConvertStringOK(s); ok {
        fmt.Println("GBK to UTF-8: ", ret)
  		return ret
    }

    return s
}
func utf8ToGBK(text string) (string, error) {
    dst := make([]byte, len(text)*2)
    tr := simplifiedchinese.GBK.NewEncoder()
    nDst, _, err := tr.Transform(dst, []byte(text), true)
    if err != nil {
        return text, err
    }
    return string(dst[:nDst]), nil
}
func GBKToUtf8(text string) (string, error) {
    dst := make([]byte, len(text)*2)
    tr := simplifiedchinese.GBK.NewDecoder()
    nDst, _, err := tr.Transform(dst, []byte(text), true)
    if err != nil {
        return text, err
    }
    return string(dst[:nDst]), nil
}
