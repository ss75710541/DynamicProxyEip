package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"runtime/debug"
	"strings"
)


// SendToMail 发送邮件
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTos := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTos, msg)
	return err
}
// SendMail 简化后的发送邮件
func SendMail(subject string, body string) {

	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	to := os.Getenv("SMTP_TO")

	fmt.Println("send email")
	err := SendToMail(user, password, host, to, subject, body, "")

	if err != nil {
		log.Println(err.Error())
		log.Println(string(debug.Stack()))
	}else{
		log.Println("Send mail success")
	}
}

//func main() {
//	subject := "test Golang to sendmail"
//	body := "11"
//	fmt.Println("send email")
//
//	SendMail(subject, body)
//
//}