package main

import (
	"fmt"
	"net/smtp"
)

func sendMail(email, token string) {
	from := "sofia.homenkova@yandex.ru"
	to := email
	subject := "Подтверждение электронной почты"
	body := fmt.Sprintf("Пожалуйста, подтвердите свою электронную почту, перейдя по следующей ссылке: http://localhost:8080/verify?token=%s", token)

	message := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body)

	err := smtp.sendMail("smtp.yandex.ru:587", smtp.PlainAuth("", from))
}
