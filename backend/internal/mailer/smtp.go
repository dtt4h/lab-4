package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type BookingReceipt struct {
	To          string
	GuestName   string
	RoomNumber  string
	RoomType    string
	CheckIn     string
	CheckOut    string
	GuestsCount int
	TotalPrice  float64
}

type BookingMailer interface {
	SendBookingReceipt(context.Context, BookingReceipt) error
}

type SMTPMailer struct{ config SMTPConfig }

func NewSMTP(config SMTPConfig) *SMTPMailer { return &SMTPMailer{config: config} }

func (m *SMTPMailer) SendBookingReceipt(ctx context.Context, receipt BookingReceipt) error {
	if m.config.Host == "" || m.config.Username == "" || m.config.Password == "" || m.config.From == "" {
		return nil // Email is optional in local development.
	}
	address, err := mail.ParseAddress(receipt.To)
	if err != nil || address.Address != receipt.To {
		return fmt.Errorf("invalid receipt recipient")
	}
	from, err := mail.ParseAddress(m.config.From)
	if err != nil {
		return fmt.Errorf("invalid SMTP sender: %w", err)
	}

	var html bytes.Buffer
	if err := receiptTemplate.Execute(&html, struct {
		BookingReceipt
		Price string
	}{receipt, formatPrice(receipt.TotalPrice)}); err != nil {
		return fmt.Errorf("render receipt: %w", err)
	}

	port := m.config.Port
	if port == 0 {
		port = 587
	}
	dialer := &net.Dialer{}
	connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(m.config.Host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("dial SMTP: %w", err)
	}
	defer connection.Close()

	client, err := smtp.NewClient(connection, m.config.Host)
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Quit()
	if ok, _ := client.Extension("STARTTLS"); !ok {
		return fmt.Errorf("SMTP server does not support STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{ServerName: m.config.Host, MinVersion: tls.VersionTLS12}); err != nil {
		return fmt.Errorf("start SMTP TLS: %w", err)
	}
	if err := client.Auth(smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)); err != nil {
		return fmt.Errorf("SMTP authentication: %w", err)
	}
	if err := client.Mail(from.Address); err != nil {
		return err
	}
	if err := client.Rcpt(address.Address); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", m.config.From, address.Address, mime.QEncoding.Encode("UTF-8", "Поиск выгодных отелей — бронирование принято"), html.String())
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}
	return writer.Close()
}

func formatPrice(price float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", price), ".00", "") + " ₽"
}

var receiptTemplate = template.Must(template.New("booking-receipt").Parse(`<!doctype html><html lang="ru"><body style="margin:0;background:#f4f7fb;font-family:Arial,sans-serif;color:#183145"><table role="presentation" width="100%" cellspacing="0" cellpadding="0"><tr><td align="center" style="padding:28px 16px"><table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="max-width:560px;background:#fff;border-radius:14px;overflow:hidden"><tr><td style="padding:28px 32px;background:#16364b;color:#fff"><h1 style="margin:0;font-size:26px">Бронирование принято</h1></td></tr><tr><td style="padding:28px 32px"><p style="margin:0 0 22px">Здравствуйте, {{.GuestName}}! Ваша заявка принята и ожидает подтверждения от отеля.</p><table role="presentation" width="100%" cellspacing="0" cellpadding="8" style="border-collapse:collapse;background:#f7fafb;border-radius:8px"><tr><td style="color:#647984">Номер</td><td align="right"><b>{{.RoomNumber}}</b>{{if .RoomType}} · {{.RoomType}}{{end}}</td></tr><tr><td style="color:#647984">Заезд</td><td align="right"><b>{{.CheckIn}}</b></td></tr><tr><td style="color:#647984">Выезд</td><td align="right"><b>{{.CheckOut}}</b></td></tr><tr><td style="color:#647984">Гостей</td><td align="right"><b>{{.GuestsCount}}</b></td></tr><tr><td style="color:#647984;border-top:1px solid #dce5e9">Итого</td><td align="right" style="border-top:1px solid #dce5e9;font-size:18px"><b>{{.Price}}</b></td></tr></table><p style="margin:24px 0 0;color:#647984;font-size:13px;line-height:1.5">Если планы изменятся, пожалуйста, свяжитесь с отелем. Это письмо сформировано автоматически.</p></td></tr></table></td></tr></table></body></html>`))
