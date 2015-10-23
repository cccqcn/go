package main
import (
        "os"
        "fmt"
	"flag"
	"bufio"
	"strconv"
	"io"
	"regexp"
	"strings"
	"io/ioutil"
)
func main() {
	flag.Parse()
	if flag.NArg() < 2 {//参数不足2个，输出错误及帮助信息
		fmt.Println("PARAMS ERROR")
		fmt.Println("V1.0")
		fmt.Println("The right params is: filename variable (offset) (precision)")
		fmt.Println("The default value of offset is 1.0")
		fmt.Println("The default value of precision is 0")
		fmt.Println("eg: file.exe a.txt VER")
		fmt.Println("eg: file.exe a.txt VER 1")
		fmt.Println("eg: file.exe a.txt VER 0.0001 4")
		return
	}
        userFile := flag.Arg(0)
	variable := flag.Arg(1)
	offset := 1.0
	if flag.NArg() > 2 {
		var offerr error
		offset, offerr = strconv.ParseFloat(flag.Arg(2), 32)
		if offerr != nil {//第三个参数增加的版本偏移解析错误
			fmt.Println(offerr)
			return
		}
	}
	precision := 0
	if flag.NArg() == 4 {//解析第四个参数精度
		precision, _ = strconv.Atoi(flag.Arg(3))
	}

	fin, err := os.OpenFile(userFile, os.O_RDWR|os.O_CREATE, 0644)
        defer fin.Close()
        if err != nil {//文件打开错误
                fmt.Println(userFile, err)
		return
        }
	fileBytes, err := ioutil.ReadFile(userFile)
	if err != nil || io.EOF == err {//读取失败
                fmt.Println(err)
		return
	}
	line := string(fileBytes)
	index := strings.Index(line, variable)
	if index == -1 {
                fmt.Println("Variable not found")
		return
	}
	preline := line[0:index]
	line = line[index:]

	wordRx := regexp.MustCompile("\"(.*?)\"")
	isFirst := false
	replacer := func(word string) string {//替换函数，将正则匹配到的内容解析为float类型并加上版本偏移
		if isFirst == true {
			return word
		}
		isFirst = true
		word2 := word[1:len(word)-1]
		oldf,_ := strconv.ParseFloat(word2, 32)
		newf := oldf + offset
		newword := "\"" + strconv.FormatFloat(newf, 'f', precision, 32) + "\""
		fmt.Print(word + "->" + newword + "\n")
		return newword
	}
	line = wordRx.ReplaceAllStringFunc(line, replacer)
	fmt.Print("\n")

	err = os.Truncate(userFile, 0) //清空文件
	if err != nil {
		fmt.Println(err)
	}
	fin.Seek(0, 0)//从头开始写入

	writer := bufio.NewWriter(fin)
	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()
	if _, err = writer.WriteString(preline + line); err != nil {//重新写入文件
                fmt.Println(err)
		return 
	}
}