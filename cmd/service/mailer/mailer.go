package mailer

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/MatthewAraujo/notify/cmd/types"
	"github.com/joho/godotenv"
)

func createEmailMessage(n types.SendEmail) string {
	message := fmt.Sprintf("Olá %s,\n\n", n.Email)
	message += fmt.Sprintf("Você recebeu uma notificação sobre o repositório '%s'.\n\n", n.RepoName)
	message += fmt.Sprintf("A ação foi realizada por %s.\n", n.Sender)
	message += fmt.Sprintf("O commit relacionado é: %s\n\n", n.Commit)
	message += "Por favor, verifique seu repositório para mais detalhes.\n\n"
	message += "Atenciosamente,\nEquipe de Notificações\n"
	return message
}
func SendMail(user types.SendEmail) {
	godotenv.Load()
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	stmpAuthor := os.Getenv("SMTP_AUTHOR")

	auth := smtp.PlainAuth(
		"",
		stmpAuthor,
		smtpPassword,
		smtpHost,
	)

	msg := createEmailMessage(user)
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		stmpAuthor,
		[]string{user.Email},
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}
