package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dchat"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type ChatService struct {
	log            *logger.Logger
	chatRepository repository.ChatRepository
	userRepository repository.UserRepository
	chatPolicy     *policy.ChatPolicy
}

func NewChatService(
	log *logger.Logger,
	chatRepository repository.ChatRepository,
	userRepository repository.UserRepository,
) *ChatService {
	return &ChatService{
		log:            log,
		chatRepository: chatRepository,
		userRepository: userRepository,
		chatPolicy:     policy.NewChatPolicy(),
	}
}

func (s *ChatService) SendMessage(sender *entity.User, req *dchat.SendMessageRequest) (*dchat.MessageResponse, error) {
	receiverID, err := id.Parse(req.ReceiverID)
	if err != nil {
		return nil, errors.ErrInvalidUUID(err)
	}

	receiver, err := s.userRepository.FindByID(receiverID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if err := s.chatPolicy.CanChat(sender, receiver); err != nil {
		return nil, err
	}

	msg := entity.NewChatMessage(sender.ID, receiverID, req.Content)

	saved, err := s.chatRepository.Save(msg)
	if err != nil {
		return nil, errors.ErrDatabase("failed to save message", err)
	}

	return toMessageResponse(saved), nil
}

func (s *ChatService) GetConversation(auth *entity.User, peerID id.UUID, limit, offset int) ([]*dchat.MessageResponse, error) {
	peer, err := s.userRepository.FindByID(peerID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if err := s.chatPolicy.CanChat(auth, peer); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	msgs, err := s.chatRepository.FindConversation(auth.ID, peerID, limit, offset)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch conversation", err)
	}

	s.chatRepository.MarkAsRead(auth.ID, peerID)

	result := make([]*dchat.MessageResponse, len(msgs))
	for i, m := range msgs {
		result[i] = toMessageResponse(m)
	}
	return result, nil
}

func toMessageResponse(m *entity.ChatMessage) *dchat.MessageResponse {
	return &dchat.MessageResponse{
		ID:         m.ID.String(),
		SenderID:   m.SenderID.String(),
		ReceiverID: m.ReceiverID.String(),
		Content:    m.Content,
		Read:       m.Read,
		CreatedAt:  m.CreatedAt.String(),
	}
}
