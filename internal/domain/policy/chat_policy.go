package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
)

const chatPolicy = "chat_policy.ChatPolicy"

type ChatPolicy struct{}

func NewChatPolicy() *ChatPolicy {
	return &ChatPolicy{}
}

func (p *ChatPolicy) CanChat(sender, receiver *entity.User) error {
	if sender.HasRole(enums.ADMIN) || receiver.HasRole(enums.ADMIN) {
		return nil
	}

	if sender.IsSeller != receiver.IsSeller {
		return nil
	}

	return apperrors.Forbidden(chatPolicy+".can_chat: users cannot chat with each other", nil)
}
