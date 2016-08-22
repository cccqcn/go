package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"ping"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Unknwon/goconfig"
)

var cfg *goconfig.ConfigFile

func main() {
	config := "config.txt"
	flag.Parse()
	if flag.NArg() == 1 {
		config = flag.Arg(0)
	}
	var err error
	cfg, err = goconfig.LoadConfigFile(config)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	fmt.Println(config + " loaded")
	sendemail := cfg.MustValue(goconfig.DEFAULT_SECTION, "sendemail")
	fmt.Println("send email when failed: " + sendemail)

	ips, err := cfg.GetSection("ips")
	if err != nil {
		log.Fatalf("GetSection err：%s", err)
	}
	fmt.Println("total ips: ", len(ips))
	fmt.Println("")
	//return
	//ip := "qq.com"
	m := make(map[string]bool)
	for _, ipportstr := range ips {

		ipport := strings.Split(ipportstr, ":")
		ip := ipport[0]
		m[ip] = true
	}
	for {
		for _, ipportstr := range ips {

			ipport := strings.Split(ipportstr, ":")
			ip := ipport[0]
			ports := []string{}
			if len(ipport) > 1 {
				ports = strings.Split(ipport[1], ",")
			}

			sendEmail := true
			if m[ip] == false {
				sendEmail = false
			}
			m[ip] = true
			fmt.Println(time.Now())
			fmt.Println("sendEmail this time: ", sendEmail)
			fmt.Println("ping " + ip)
			pingflag := ping.Ping(ip, 5)
			fmt.Println("result: " + strconv.FormatBool(pingflag))

			if pingflag == false {
				if sendEmail == true {
					email(true, ip, "")
				}
				m[ip] = false
			} else if len(ports) > 0 {

				for _, port := range ports {
					ipwithport := ip + ":" + port
					fmt.Println("DialTimeout " + ipwithport)
					portflag := PortIsOpen(ipwithport, 3)
					fmt.Println("PortIsOpen: " + strconv.FormatBool(portflag))

					if portflag == false {
						if sendEmail == true {
							email(false, ip, port)
						}
						m[ip] = false
					}
				}
			}
		}

		intervalInSeconds := cfg.MustInt(goconfig.DEFAULT_SECTION, "intervalInSeconds")
		fmt.Println("waiting for: " + strconv.Itoa(intervalInSeconds) + "s")
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)

		fmt.Println("")
	}
	//fmt.Println("Telnet " + ipport)
	//buf, err := Telnet([]string{"w_Hello World", "r_50", "r_30", "r_30"}, ipport, 5)
	//fmt.Println(err)
	//fmt.Println(string(buf))
	//socketTest(ipport)
}

//This function is currently not used.
func socketTest(ip string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("connect success")
	words := "hello world!"
	conn.Write([]byte(words))
	fmt.Println("send over")
}

func PortIsOpen(ip string, timeout int) bool {
	//con, err := net.Dial("tcp", ip)
	con, err := net.DialTimeout("tcp", ip, time.Duration(timeout)*time.Second)
	if err != nil {
		fmt.Println(err)
		return false
	}
	con.Close()
	return true
}

//This function is currently not used.
func Telnet(action []string, ip string, timeout int) (buf []byte, err error) {
	con, err := net.Dial("tcp", ip)
	//con, err := net.DialTimeout("tcp", ip, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("111111")
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(time.Second * 5))
	for _, v := range action {
		l := strings.SplitN(v, "_", 2)
		if len(l) < 2 {
			return
		}
		switch l[0] {
		case "r":
			var n int
			n, err = strconv.Atoi(l[1])
			if err != nil {
				return
			}
			p := make([]byte, n)
			n, err = con.Read(p)
			if err != nil {
				return
			}
			buf = append(buf, p[:n]...)
			fmt.Println(buf)
		case "w":
			_, err = con.Write([]byte(l[1]))
		}
	}
	return
}

func email(isPing bool, ip string, port string) {
	sendemail := cfg.MustBool(goconfig.DEFAULT_SECTION, "sendemail")
	if sendemail == false {
		return
	}
	from := cfg.MustValue(goconfig.DEFAULT_SECTION, "from")
	to := cfg.MustValue(goconfig.DEFAULT_SECTION, "to")
	toEmails := strings.Split(to, ",")
	smtpport := cfg.MustValue(goconfig.DEFAULT_SECTION, "smtp")
	smtpport2 := strings.Split(smtpport, ":")
	smtpip := smtpport2[0]
	username := cfg.MustValue(goconfig.DEFAULT_SECTION, "username")
	password := cfg.MustValue(goconfig.DEFAULT_SECTION, "password")

	subject := ""
	if isPing == true {
		subject = "服务器" + ip + "无法ping通"
	} else {
		subject = "服务器" + ip + "端口" + port + "无法连接"
	}

	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		from,
		strings.Join(toEmails, ","),
		subject,
		"见邮件主题.",
	}

	buffer := new(bytes.Buffer)
	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		smtpip,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpport,
		auth,
		from,
		toEmails,
		buffer.Bytes(),
	)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("send OK")
}

func emailScript() (script string) {
	return "From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\n\r\n{{.Message}}"
}
