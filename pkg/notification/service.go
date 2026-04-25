package notification

import (
	"encoding/json"
	"fmt"
)

type NotificationService interface {
	SendTicketEmail(toEmail string, eventID string, seatID string) error
	SendTicketsEmail(toEmail string, tickets []TicketData) error
}

type TicketData struct {
	EventID    string
	SeatID     string
	SeatNumber string
}

type notificationService struct {
	// We can add Gmail API client here later
}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) SendTicketEmail(toEmail string, eventID string, seatID string) error {
	// 1. Generate QR content in JSON format
	qrData := map[string]string{
		"eventId": eventID,
		"seatId":  seatID,
	}
	qrJson, err := json.Marshal(qrData)
	if err != nil {
		return fmt.Errorf("failed to marshal QR data: %v", err)
	}

	// 2. Generate QR bytes
	qrBytes, err := GenerateQRCode(string(qrJson))
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	// 3. Send via SMTP
	return s.sendEmail(toEmail, "Your Event Ticket", "Please find your ticket QR code attached.", qrBytes)
}
func (s *notificationService) SendTicketsEmail(toEmail string, tickets []TicketData) error {
	var attachments [][]byte
	var filenames []string

	for i, t := range tickets {
		qrData := map[string]string{
			"eventId": t.EventID,
			"seatId":  t.SeatID,
		}
		qrJson, _ := json.Marshal(qrData)
		qrBytes, err := GenerateQRCode(string(qrJson))
		if err != nil {
			return err
		}
		attachments = append(attachments, qrBytes)
		filenames = append(filenames, fmt.Sprintf("ticket_%s.png", t.SeatNumber))
		_ = i
	}

	body := fmt.Sprintf("You have successfully booked %d seats. Please find your tickets attached.", len(tickets))
	return s.sendEmailMulti(toEmail, "Your Event Tickets", body, attachments, filenames)
}
