package models

type ChatUser struct {
	UserId          string
	Nickname        string
	LastMessageTime *string // nil if no messages exchanged, *string (pointer to distinguish null/no-messages from empty string)
}