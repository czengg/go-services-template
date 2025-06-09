package upwardli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	webhookSDK "elevate/packages/webhook-go"

	"go.uber.org/zap"
)

func writeErrorResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	errResp, _ := json.Marshal(map[string]string{
		"error": errorMessage,
	})
	w.Write(errResp)
}

func WithWebhookVerification(verifier webhookSDK.WebhookVerifier, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				writeErrorResponse(w, "error reading body")
				return
			}

			signature := r.Header.Get("Upwardli-Signature")
			split := strings.Split(signature, ",")

			logger.Debug("upwardli webhook called")

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			mapping := make(map[string]string)
			for _, element := range split {
				keyValue := strings.Split(element, "=")
				if len(keyValue) != 2 {
					writeErrorResponse(w, "invalid signature format")
					return
				}
				mapping[keyValue[0]] = keyValue[1]
			}

			requiredFields := []string{"t", "v1"}
			for _, field := range requiredFields {
				if val, ok := mapping[field]; !ok || val == "" {
					writeErrorResponse(w, fmt.Sprintf("missing required field: %s", field))
					return
				}
			}

			compare := mapping["t"] + "." + string(body)
			if !verifier.VerifyWebhook([]byte(compare), []byte(mapping["v1"])) {
				writeErrorResponse(w, "invalid signature")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
