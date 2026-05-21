package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const chatDB = "chat_db.ChatDB"

type ChatDB struct {
	db *gorm.DB
}

func NewChatDB(db *gorm.DB) repository.ChatRepository {
	return &ChatDB{db: db}
}

func (c *ChatDB) Save(msg *entity.ChatMessage) (*entity.ChatMessage, error) {
	if err := c.db.WithContext(context.Background()).Create(msg).Error; err != nil {
		return nil, apperrors.HandlePgError(chatDB+".save", err)
	}
	return msg, nil
}

func (c *ChatDB) FindConversation(userA, userB id.UUID, limit, offset int) ([]*entity.ChatMessage, error) {
	var msgs []*entity.ChatMessage
	err := c.db.WithContext(context.Background()).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userA, userB, userB, userA).
		Order("created_at asc").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	if err != nil {
		return nil, apperrors.HandlePgError(chatDB+".find_conversation", err)
	}
	return msgs, nil
}

func (c *ChatDB) MarkAsRead(receiverID, senderID id.UUID) error {
	return c.db.WithContext(context.Background()).Model(&entity.ChatMessage{}).
		Where("receiver_id = ? AND sender_id = ? AND read = false", receiverID, senderID).
		Update("read", true).Error
}

func (c *ChatDB) UnreadCount(receiverID id.UUID) (int64, error) {
	var count int64
	err := c.db.WithContext(context.Background()).Model(&entity.ChatMessage{}).
		Where("receiver_id = ? AND read = false", receiverID).
		Count(&count).Error
	return count, err
}
