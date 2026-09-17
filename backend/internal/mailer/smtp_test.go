package mailer

import (
	"bytes"
	"strings"
	"testing"
)

func TestReceiptTemplateEscapesGuestData(t *testing.T) {
	var body bytes.Buffer
	err := receiptTemplate.Execute(&body, struct {
		BookingReceipt
		Price string
	}{BookingReceipt{GuestName: "<script>alert(1)</script>", RoomNumber: "203", CheckIn: "2026-10-01", CheckOut: "2026-10-03", GuestsCount: 2}, "7 000 ₽"})
	if err != nil {
		t.Fatalf("render receipt: %v", err)
	}
	if strings.Contains(body.String(), "<script>") {
		t.Fatal("guest data must be HTML-escaped")
	}
	if !strings.Contains(body.String(), "203") {
		t.Fatal("receipt must include room number")
	}
}
