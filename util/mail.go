package util

import "gopkg.in/gomail.v2"

const Email="QEHkTGtfDWU2Y9J@gmail.com"
const Username="QEHkTGtfDWU2Y9J"
const Password="LCpPRP7V^4GYjWkn"
func SendVerifyMail(email string,header string, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From",email )
	m.SetHeader("To", email)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", header)
	m.SetBody("text/html", content)
	d := gomail.NewDialer("smtp.gmail.com", 465,Username , Password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
