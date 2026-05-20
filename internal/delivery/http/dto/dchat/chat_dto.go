package dchat

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

func (r *SendMessageRequest) Validate() error {
	if r.ReceiverID == "" {
		return errors.ErrBadRequest("receiver_id is required", nil)
	}
	if len(r.Content) == 0 {
		return errors.ErrBadRequest("content is required", nil)
	}
	if len(r.Content) > 4000 {
		return errors.ErrBadRequest("content must not exceed 4000 characters", nil)
	}
	return nil
}

type MessageResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Read       bool   `json:"read"`
	CreatedAt  string `json:"created_at"`
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload MessageResponse `json:"payload"`
}

type chatError struct{ msg string }

func (e *chatError) Error() string { return e.msg }
