package dchat

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

func (r *SendMessageRequest) Validate() error {
	if r.ReceiverID == "" {
		return errBadRequest("receiver_id is required")
	}
	if len(r.Content) == 0 {
		return errBadRequest("content is required")
	}
	if len(r.Content) > 4000 {
		return errBadRequest("content must not exceed 4000 characters")
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

func errBadRequest(msg string) error {
	return &chatError{msg: msg}
}

type chatError struct{ msg string }

func (e *chatError) Error() string { return e.msg }
