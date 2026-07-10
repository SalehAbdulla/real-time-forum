package models

type Notification struct {
	NotificationId int
	UserId         string
	ActorId        string
	ActorNickname  string
	EntityType     string // "comment" or "message"
	EntityId       int
	IsRead         int // 0 or 1
	CreatedAt      string
}