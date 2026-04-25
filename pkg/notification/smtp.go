package notification

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

func (s *notificationService) sendEmail(to, subject, body string, attachment []byte) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD") // นี่คือ App Password 16 หลัก
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	if from == "" || password == "" {
		fmt.Printf("Warning: SMTP_EMAIL or SMTP_PASSWORD not set. Skipping email.\n")
		fmt.Printf("Simulation: Sending ticket to %s with %d bytes attachment\n", to, len(attachment))
		return nil
	}

	// สร้าง Message พร้อม Attachment (ใช้รูปแบบเดิมที่เคยเขียนไว้)
	boundary := "my-boundary-777"
	messageBody := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: image/png\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"Content-Disposition: attachment; filename=\"ticket.png\"\r\n\r\n"+
			"%s\r\n"+
			"--%s--",
		from, to, subject, boundary, boundary, body, boundary, base64.StdEncoding.EncodeToString(attachment), boundary,
	)

	// Authentication
	fmt.Printf("[SMTP] Authenticating with %s...\n", from)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมล
	fmt.Printf("[SMTP] Sending email to %s via %s:%s...\n", to, smtpHost, smtpPort)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(messageBody))
	if err != nil {
		return fmt.Errorf("failed to send smtp email: %v", err)
	}

	fmt.Printf("[SMTP] Success! Email sent to %s\n", to)
	return nil
}
func (s *notificationService) sendEmailMulti(to, subject, body string, attachments [][]byte, filenames []string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	if from == "" || password == "" {
		fmt.Printf("Warning: SMTP_EMAIL or SMTP_PASSWORD not set. Skipping email.\n")
		return nil
	}

	boundary := "my-boundary-999"
	header := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
			"%s\r\n\r\n",
		from, to, subject, boundary, boundary, body,
	)

	var msg bytes.Buffer
	msg.WriteString(header)

	for i, attachment := range attachments {
		filename := "ticket.png"
		if i < len(filenames) {
			filename = filenames[i]
		}
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: image/png\r\n")
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename))
		msg.WriteString(base64.StdEncoding.EncodeToString(attachment))
		msg.WriteString("\r\n")
	}
	msg.WriteString(fmt.Sprintf("--%s--", boundary))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send multi-attachment email: %v", err)
	}

	fmt.Printf("[SMTP] Success! Multi-ticket email sent to %s with %d attachments\n", to, len(attachments))
	return nil
}
