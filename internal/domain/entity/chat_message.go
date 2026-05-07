package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ChatMessage struct {
	ID         id.UUID   `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index"`
	SenderID   id.UUID   `gorm:"index;not null"`
	ReceiverID id.UUID   `gorm:"index;not null"`
	Content    string    `gorm:"size:4000;not null"`
	Read       bool      `gorm:"default:false"`
}

func NewChatMessage(senderID, receiverID id.UUID, content string) *ChatMessage {
	return &ChatMessage{
		ID:         id.NewUUID(),
		CreatedAt:  time.Now(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Read:       false,
	}
}
