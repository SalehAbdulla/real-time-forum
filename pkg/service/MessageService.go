package service

import (
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/payload/message"
	db "real-time-forum/pkg/repositories"
)

type MessageService interface {
	GetChatUsers(currentUserID string) ([]message.ChatUserDTO, error)
	GetMessages(conversationPartnerID string, currentUserID string, offset int, limit int) (message.MessagesResponse, error)
}

type MessageServiceImpl struct {
	messageRepo    db.MessageRepository
	authRepository db.AuthRepository
	sessionManager *SessionManager
}

func NewMessageService(messageRepo db.MessageRepository, authRepository db.AuthRepository) MessageService {
	return MessageServiceImpl{
		messageRepo:    messageRepo,
		authRepository: authRepository,
		sessionManager: DefaultSessionManager,
	}
}

func (m MessageServiceImpl) GetMessages(conversationPartnerID string, currentUserID string, offset int, limit int) (message.MessagesResponse, error) {
	if conversationPartnerID == currentUserID {
		return message.MessagesResponse{}, realtimeforum.ErrBadRequest
	}

	if err := m.authRepository.DoesUserExists(conversationPartnerID); err != nil {
		return message.MessagesResponse{}, err
	}

	messages, totalElements, err := m.messageRepo.GetMessages(conversationPartnerID, currentUserID, offset, limit)
	if err != nil {
		return message.MessagesResponse{}, err
	}

	dtos := make([]message.MessageDTO, len(messages))
	for i, msg := range messages {
		dtos[i] = message.MessageDTO{
			MessageId:   msg.MessageId,
			SenderId:    msg.SenderId,
			RecipientId: msg.RecipientId,
			TextMessage: msg.TextMessage,
			TimeStamp:   msg.TimeStamp,
			IsRead:      msg.IsRead,
		}
	}

	return message.MessagesResponse{
		Messages:      dtos,
		Offset:        offset,
		Limit:         limit,
		TotalElements: totalElements,
	}, nil
}

func (m MessageServiceImpl) GetChatUsers(currentUserID string) ([]message.ChatUserDTO, error) {
	chatUsers, err := m.messageRepo.GetChatUsers(currentUserID)
	if err != nil {
		return nil, err
	}

	dtos := make([]message.ChatUserDTO, len(chatUsers))
	for i, cu := range chatUsers {
		isOnline := 0
		if m.sessionManager.IsUserOnline(cu.UserId) {
			isOnline = 1
		}

		lastMessageTime := ""
		if cu.LastMessageTime != nil {
			lastMessageTime = *cu.LastMessageTime
		}

		dtos[i] = message.ChatUserDTO{
			UserId:          cu.UserId,
			Nickname:        cu.Nickname,
			IsOnline:        isOnline,
			LastMessageTime: lastMessageTime,
		}
	}

	return dtos, nil
}
