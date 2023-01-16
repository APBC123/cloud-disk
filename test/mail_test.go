package test

import (
	"cloud-disk/core/define"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendMail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <2290312980@qq.com>"
	e.To = []string{"2290312980@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>123456</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "2290312980@qq.com", define.MailPassword, "smtp.qq.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	//InsecureSkipVerify: true 用于跳过安全验证
	if err != nil {
		t.Fatal(err)
	}
}
