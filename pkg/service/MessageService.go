package service

import (
	"real-time-forum/pkg/payload/message"
	db "real-time-forum/pkg/repositories"
)

type MessageService interface {
	GetChatUsers(currentUserID string) ([]message.ChatUserDTO, error)
}

type MessageServiceImpl struct {
	messageRepo    db.MessageRepository
	sessionManager *SessionManager
}

func NewMessageService(messageRepo db.MessageRepository) MessageService {
	return MessageServiceImpl{
		messageRepo:    messageRepo,
		sessionManager: DefaultSessionManager,
	}
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