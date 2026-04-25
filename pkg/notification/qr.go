package notification

import (
	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code from the given content and returns the PNG bytes
func GenerateQRCode(content string) ([]byte, error) {
	return qrcode.Encode(content, qrcode.Medium, 256)
}
