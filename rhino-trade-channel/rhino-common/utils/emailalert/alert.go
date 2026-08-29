package emailalert

import (
	"crypto/tls"
	"net"
	"rhino-common/utils/cache"
	"rhino-common/utils/envutil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	smtp_host string = envutil.GetEnvVarValue("SMTP_HOST","10.2.75.38")
	smtp_port string = envutil.GetEnvVarValue("SMTP_PORT","2525")
	smtp_user string = envutil.GetEnvVarValue("SMTP_USER","eagle")
	smtp_pwd  string = envutil.GetEnvVarValue("SMTP_PWD","gf37888676")
	smtp_from string = envutil.GetEnvVarValue("SMTP_FROM","eagle@gf.com.cn")
	smtp_to   string = envutil.GetEnvVarValue("SMTP_TO","linchunquan@gf.com.cn")
	smtp_conn *gomail.Dialer
	alertCache *cache.DurationCache
	hostIp string
)


func init(){
	port, _ := strconv.Atoi(smtp_port)
	smtp_conn = gomail.NewDialer(smtp_host, port, smtp_user, smtp_pwd)
	smtp_conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	alertCache = cache.NewDurationCache(2*time.Hour, func(key string,value interface{}){})

	getHostIp()
}

func Send(subject, content string) error{
	var err error
	_, cached := alertCache.Get(subject)
	if !cached{
		m := gomail.NewMessage()
		m.SetHeader("From", smtp_from)
		receivers := strings.Split(smtp_to, ",")
		m.SetHeader("To",   receivers...)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", "hostIP:"+hostIp+" ===> \n"+content)
		err = smtp_conn.DialAndSend(m)
		if err==nil{
			alertCache.Put(subject, "")
		}
	}
	return err
}

func SendTo(subject, content string, to []string) error{
	var err error
	_, cached := alertCache.Get(subject)
	if !cached{
		m := gomail.NewMessage()
		m.SetHeader("From", smtp_from)
		m.SetHeader("To",   to...)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", "hostIP:"+hostIp+" ===> "+content)
		err = smtp_conn.DialAndSend(m)
		if err==nil{
			alertCache.Put(subject, "")
		}
	}
	return err
}

func SendEmail(subject, content string, to []string, cc[]string, atts[]string) error{
	m := gomail.NewMessage()
	m.SetHeader("From", smtp_from)
	m.SetHeader("To",   to...)
	m.SetHeader("Cc",   cc...)
	m.SetHeader("Subject", subject)
	for _, att := range atts {
		m.Attach(att)
	}
	//m.SetBody("text/plain", "hostIP:"+hostIp+" ===> "+content)
	if strings.Contains(content, "<span") {
		m.SetBody("text/html", content)
	} else {
		m.SetBody("text/plain", content)
	}
	err := smtp_conn.DialAndSend(m)
	return err
}

func getHostIp(){
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		var ipSet []string
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ipSet = append(ipSet, ipnet.IP.String())
				}
			}
		}
		hostIp = strings.Join(ipSet,",")
	}
}

func GetEagleUser() string {
	return smtp_user
}

func GetEagleUserPwd() string {
	return smtp_pwd
}