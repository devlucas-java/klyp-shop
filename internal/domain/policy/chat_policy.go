package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

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

	return errors.ErrForbidden(nil)
}
