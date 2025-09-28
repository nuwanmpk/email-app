// internal/email/sender.go
package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// SendEmail sends a plain-text email using SMTP.
// It supports implicit TLS on port 465 and STARTTLS on port 587.
func SendEmail(subject, emailBody string) error {
	from := os.Getenv("SMTP_EMAIL")    // e.g. your_email@gmail.com (or SMTP username)
	password := os.Getenv("SMTP_PASS") // App password or SMTP password
	to := os.Getenv("EMAIL_TO")        // recipient email
	smtpHost := os.Getenv("SMTP_HOST") // e.g. smtp.gmail.com
	smtpPort := os.Getenv("SMTP_PORT") // "465" for implicit TLS or "587" for STARTTLS

	if from == "" || password == "" || to == "" || smtpHost == "" || smtpPort == "" {
		return errors.New("missing required SMTP env vars: SMTP_EMAIL, SMTP_PASS, EMAIL_TO, SMTP_HOST, SMTP_PORT")
	}

	// Normalize port
	port := strings.TrimSpace(smtpPort)

	// Build RFC 5322 style message (simple plain text)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		to, subject, emailBody,
	))

	// Common SMTP auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Choose flow by port
	switch port {
	case "465":
		// Implicit TLS from the start
		tlsCfg := &tls.Config{
			ServerName:         smtpHost,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false, // keep secure; only set true for local/self-signed testing
		}

		conn, err := tls.Dial("tcp", net.JoinHostPort(smtpHost, port), tlsCfg)
		if err != nil {
			return fmt.Errorf("tls dial failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			return fmt.Errorf("smtp new client failed: %w", err)
		}
		defer client.Quit()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
		if err = client.Mail(from); err != nil {
			return fmt.Errorf("smtp MAIL FROM failed: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp RCPT TO failed: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp DATA open failed: %w", err)
		}
		if _, err = w.Write(msg); err != nil {
			_ = w.Close()
			return fmt.Errorf("smtp DATA write failed: %w", err)
		}
		if err = w.Close(); err != nil {
			return fmt.Errorf("smtp DATA close failed: %w", err)
		}
		return nil

	case "587":
		// Plain TCP, then STARTTLS
		dialer := &net.Dialer{Timeout: 10 * time.Second}
		conn, err := dialer.Dial("tcp", net.JoinHostPort(smtpHost, port))
		if err != nil {
			return fmt.Errorf("tcp dial failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			return fmt.Errorf("smtp new client failed: %w", err)
		}
		defer client.Quit()

		// Upgrade to TLS
		tlsCfg := &tls.Config{
			ServerName: smtpHost,
			MinVersion: tls.VersionTLS12,
		}
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("server does not support STARTTLS on port 587")
		}
		if err = client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
		if err = client.Mail(from); err != nil {
			return fmt.Errorf("smtp MAIL FROM failed: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp RCPT TO failed: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp DATA open failed: %w", err)
		}
		if _, err = w.Write(msg); err != nil {
			_ = w.Close()
			return fmt.Errorf("smtp DATA write failed: %w", err)
		}
		if err = w.Close(); err != nil {
			return fmt.Errorf("smtp DATA close failed: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("unsupported SMTP_PORT %q (use 465 or 587)", port)
	}
}
