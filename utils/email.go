package utils

import (
	"CloudMusic/config"
	"CloudMusic/global"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
)

var htmlBody = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>欢迎注册 Se7enMusic</title>
  <style>
    body{margin:0;padding:0;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,"Noto Sans",sans-serif;color:#333;background:#f7f7f7;}
    .wrap{max-width:600px;margin:40px auto;padding:30px 20px;background:#fff;border-radius:8px;box-shadow:0 4px 12px rgba(0,0,0,.08);}
    .header{text-align:center;font-size:22px;font-weight:bold;color:#222;margin-bottom:20px;}
    .code{font-size:28px;font-weight:bold;color:#ff3b30;letter-spacing:4px;background:#fafafa;border:1px dashed #ff3b30;border-radius:4px;padding:12px 0;text-align:center;margin:20px 0;}
    .desc{font-size:14px;color:#666;line-height:1.6;margin-bottom:25px;}
    .footer{font-size:12px;color:#aaa;text-align:center;margin-top:25px;}
  </style>
</head>
<body>
  <div class="wrap">
    <div class="header">欢迎注册 Se7enMusic</div>
    <p class="desc">	
      你好，%s！<br>
      感谢注册 Se7enMusic，请在 5 分钟内输入下面的验证码完成注册：
    </p>
    <div class="code">%s</div>
    <p class="desc">
      如非本人操作，请忽略此邮件。<br>
    </p>
    <div class="footer">Se7enMusic © 2025</div>
  </div>
</body>
</html>`

func InitD() {
	mailConfig := config.AppConfig.Email
	d := gomail.NewDialer(mailConfig.MailHost, mailConfig.MailPort, mailConfig.MailUser, mailConfig.MailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	global.D = d
}

func SendCodeEmail(toMail string, code string) error {
	body := fmt.Sprintf(htmlBody, toMail, code)
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.Email.MailFrom)
	m.SetHeader("To", toMail)
	m.SetHeader("Subject", "Se7enMusic——注册")
	m.SetBody("text/html", body)
	err := global.D.DialAndSend(m)
	if err != nil {
		log.Printf("Send Mail err===>%v", err)
		return err
	}
	return nil
}
