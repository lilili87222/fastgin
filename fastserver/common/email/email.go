package email

import (
	"errors"
	"github.com/spf13/viper"
	mail "github.com/xhit/go-simple-mail/v2"
	"strings"
	"text/template"
)

func getSmtpClient() (*mail.SMTPClient, error) {
	//if smtpClient == nil {
	server := mail.NewSMTPClient()
	server.Host = viper.GetString("email.host")
	server.Port = viper.GetInt("email.port")
	server.Username = viper.GetString("email.account")
	server.Password = viper.GetString("email.password")
	if viper.GetBool("email.is-ssl") {
		server.Encryption = mail.EncryptionSTARTTLS
	}
	client, e := server.Connect()
	return client, e
}

// SendAccountEmail 用户注册信息发送到用户邮箱上
func SendAccountEmail(userEmail, subject string, content string) error {
	blackEmails := viper.GetStringSlice("email.blackmail")
	for _, email := range blackEmails {
		if strings.HasSuffix(userEmail, email) {
			return errors.New("email address is blocked, use other email type please")
		}
	}
	from := viper.GetString("email.from")

	// Create email
	email := mail.NewMSG()
	email.SetFrom(from)
	email.AddTo(userEmail)
	//email.AddCc("129TO139Test")
	email.SetSubject(subject)
	email.SetBody(mail.TextHTML, content) //发送html信息
	//email.AddAttachment("super_cool_file.png") // 附件
	client, e := getSmtpClient()
	if e != nil {
		return e
	}
	defer client.Close()
	return email.Send(client)
}
func NewEmail(templateName string, data map[string]any) (string, error) {
	tfileName := templateName + ".email.html"
	tpl, err := template.New(tfileName).ParseFiles("conf/" + tfileName)
	if err != nil {
		return "", err
	}
	var body strings.Builder
	err = tpl.Execute(&body, data)
	return body.String(), err
}

func SendRegisterEmail(toEmail, subject, code string) error {
	content, err := NewEmail("register", map[string]any{"code": code, "email": toEmail})
	if err != nil {
		return err
	}
	return SendAccountEmail(toEmail, subject, content)
}
