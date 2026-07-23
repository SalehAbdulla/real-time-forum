package websocket

import "encoding/json"

const (
	MsgTypePrivateMsg = "private_msg"

	MsgTypeIncomingMsg  = "incoming_msg"
	MsgTypeUserStatus   = "user_status"
	MsgTypeNotification = "notification"
	MsgTypeSendError    = "send_error"
)

type WSMessage struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	SenderId  string          `json:"senderId,omitempty"`
	RecipientId string        `json:"recipientId,omitempty"`
	Text       string         `json:"text,omitempty"`
	TimeStamp  string         `json:"timeStamp,omitempty"`
	SenderNickname string     `json:"senderNickname,omitempty"`
}

type PrivateMsgPayload struct {
	RecipientId string `json:"recipientId"`
	Text        string `json:"text"`
}

type IncomingMsgPayload struct {
	MessageId       int    `json:"messageId"`
	SenderId        string `json:"senderId"`
	SenderNickname  string `json:"senderNickname"`
	Text            string `json:"text"`
	TimeStamp       string `json:"timeStamp"`
}

type UserStatusPayload struct {
	UserId   string `json:"userId"`
	IsOnline int    `json:"isOnline"` // 1 or 0
}

type NotificationPayload struct {
	NotificationId  int    `json:"notificationId"`
	ActorId         string `json:"actorId"`
	ActorNickname   string `json:"actorNickname"`
	EntityType      string `json:"entityType"` // "comment" or "message"
	EntityId        int    `json:"entityId"`
	CreatedAt       string `json:"createdAt"`
}

type SendErrorPayload struct {
	RecipientId string `json:"recipientId"`
	Message     string `json:"message"`
}
