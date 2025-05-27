package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/upstash/qstash-go"
)

type QStashClient struct {
	client   *qstash.Client
	receiver *qstash.Receiver
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewQStashClient(token, signingKey, nextSigningKey string) (*QStashClient, error) {
	// Initialize client
	client := qstash.NewClient(token)

	// Initialize receiver
	receiver := qstash.NewReceiver(signingKey, nextSigningKey)

	return &QStashClient{
		client:   client,
		receiver: receiver,
	}, nil
}

// PublishMessage sends a message to QStash
func (q *QStashClient) PublishMessage(ctx context.Context, url string, msgType string, payload interface{}) error {
	msg := map[string]interface{}{
		"type":    msgType,
		"payload": payload,
	}

	res, err := q.client.PublishJSON(qstash.PublishJSONOptions{
		Url:  url,
		Body: msg,
	})
	if err != nil {
		return err
	}
	_ = res // You can use res.MessageId if needed
	return nil
}

// PublishDelayedMessage sends a message to QStash with a delay
func (q *QStashClient) PublishDelayedMessage(ctx context.Context, url string, msgType string, payload interface{}, delay time.Duration) error {
	msg := map[string]interface{}{
		"type":    msgType,
		"payload": payload,
	}

	// Convert delay to ISO 8601 duration string (e.g., "PT1H" for 1 hour)
	delayStr := formatDuration(delay)

	res, err := q.client.PublishJSON(qstash.PublishJSONOptions{
		Url:   url,
		Body:  msg,
		Delay: delayStr,
	})
	if err != nil {
		return err
	}
	_ = res // You can use res.MessageId if needed
	return nil
}

// VerifyRequest verifies an incoming QStash request
func (q *QStashClient) VerifyRequest(signature string, body []byte) error {
	return q.receiver.Verify(qstash.VerifyOptions{
		Signature: signature,
		Body:      string(body),
	})
}

// HandleRequest processes an incoming QStash request
func (q *QStashClient) HandleRequest(handler func(context.Context, *Message)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get signature from header
		signature := r.Header.Get("Upstash-Signature")
		if signature == "" {
			http.Error(w, "Missing signature", http.StatusBadRequest)
			return
		}

		// Read body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Verify request
		err = q.VerifyRequest(signature, bodyBytes)
		if err != nil {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(bodyBytes, &msg); err != nil {
			http.Error(w, "Invalid message format", http.StatusBadRequest)
			return
		}

		// Handle message
		handler(r.Context(), &msg)

		// Return success
		w.WriteHeader(http.StatusOK)
	}
}

// formatDuration converts a time.Duration to an ISO 8601 duration string
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("PT%dS", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("PT%dM", minutes)
	}
	hours := minutes / 60
	return fmt.Sprintf("PT%dH", hours)
}
