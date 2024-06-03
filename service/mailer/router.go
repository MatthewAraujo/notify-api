package mailer

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/types"
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

func setupEmail() (smtp.Auth, string) {
	godotenv.Load()
	smtpPassword := config.Envs.SMTP.Password
	stmpAuthor := config.Envs.SMTP.Author
	smtpHost := config.Envs.SMTP.Host
	smtpPort := config.Envs.SMTP.Port

	auth := smtp.PlainAuth(
		"",
		stmpAuthor,
		smtpPassword,
		smtpHost,
	)

	return auth, smtpPort

}

func sendEmail(msg string, user types.SendEmail, smtpHost, smtpPort, stmpAuthor string, auth smtp.Auth) {
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

func SendWelcomeEmail(user types.WelcomeEmail) {
	msg := fmt.Sprintf("Olá %s,\n\n", user.Owner)
	msg += fmt.Sprintf("Seja bem-vindo ao repositório %s.\n\n", user.Repository)
	msg += "Agora você receberá notificações sobre as ações realizadas no repositório.\n\n"
	msg += "Atenciosamente,\nEquipe de Notificações\n"

	auth, _ := setupEmail()
	sendEmail(msg, types.SendEmail{Email: user.Email}, config.Envs.SMTP.Host, config.Envs.SMTP.Port, config.Envs.SMTP.Author, auth)
	log.Println("Email sent to", user.Email)
}

func SendMail(user types.SendEmail) {
	msg := createEmailMessage(user)

	auth, _ := setupEmail()
	sendEmail(msg, user, config.Envs.SMTP.Host, config.Envs.SMTP.Port, config.Envs.SMTP.Author, auth)
	log.Println("Email sent to", user.Email)
}
