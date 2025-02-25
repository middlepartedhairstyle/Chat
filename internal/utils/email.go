package utils

import "github.com/go-mail/mail"

// EmailSendCode 邮箱发送验证码
func EmailSendCode(email string, code string) (bool, error) {
	m := mail.NewMessage()
	m.SetHeader("From", Cfg.Email.From)
	m.SetHeader("To", email)
	//m.SetAddressHeader("Cc", "2731337079@qq.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "验证码为"+code)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := mail.NewDialer(Cfg.Email.Host, Cfg.Email.Port, Cfg.Email.From, Cfg.Email.Password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
	return true, nil
}
