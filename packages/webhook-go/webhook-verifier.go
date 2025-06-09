package webhook_go

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
)

type WebhookVerifier interface {
	VerifyWebhook(message, signature []byte) bool
}

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
	}
}

func (v *Verifier) VerifyWebhook(body, signature []byte) bool {
	mac := hmac.New(crypto.SHA256.New, []byte(v.secret))
	mac.Write(body)

	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal(signature, []byte(expectedMAC))
}
