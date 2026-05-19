package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

// ChatPolicy contém as regras de negócio para o chat entre usuários.
type ChatPolicy struct{}

func NewChatPolicy() *ChatPolicy {
	return &ChatPolicy{}
}

// CanChat verifica se o remetente pode enviar mensagem ao destinatário.
// Regras:
//   - Admin pode conversar com qualquer um.
//   - Seller pode conversar com buyer (e vice-versa).
//   - Seller não pode conversar com outro seller.
//   - Buyer não pode conversar com outro buyer.
func (p *ChatPolicy) CanChat(sender, receiver *entity.User) error {
	if sender.HasRole(enums.ADMIN) || receiver.HasRole(enums.ADMIN) {
		return nil
	}

	// sellers falam com buyers e vice-versa
	if sender.IsSeller != receiver.IsSeller {
		return nil
	}

	return errors.ErrForbidden(nil)
}
