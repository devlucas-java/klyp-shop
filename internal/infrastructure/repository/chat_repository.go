package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ChatRepository interface {
	Save(msg *entity.ChatMessage) (*entity.ChatMessage, error)
	FindConversation(userA, userB id.UUID, limit, offset int) ([]*entity.ChatMessage, error)
	MarkAsRead(receiverID, senderID id.UUID) error
	UnreadCount(receiverID id.UUID) (int64, error)
}
