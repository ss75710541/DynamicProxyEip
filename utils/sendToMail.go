package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func SendMail(subject string, body string) {

	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	to := os.Getenv("SMTP_TO")

	fmt.Println("send email")
	err := SendToMail(user, password, host, to, subject, body, "")

	if err != nil {
		log.Panic(err)
	}else{
		log.Println("Send mail success")
	}
}

//func main() {
//	to := "liujinye@hisuntech.com"
//	subject := "test Golang to sendmail"
//	body := "11"
//	fmt.Println("send email")
//	err := SendToMail(user, password, host, to, subject, body, "")
//
//	if err != nil {
//		log.Panic(err)
//	}else{
//		log.Println("Send mail success")
//	}
//}