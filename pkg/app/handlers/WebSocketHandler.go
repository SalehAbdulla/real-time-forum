package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"real-time-forum/pkg/app/service"
	pkgwebsocket "real-time-forum/pkg/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func (re *Repository) ServeWs(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("session_token")
	if err != nil || token.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := service.DefaultSessionManager.GetUserIdByToken(token.Value)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}


	client := &pkgwebsocket.Client{
		Hub:    re.Hub,
		Conn:   conn,
		Send:   make(chan []byte, pkgwebsocket.SendBufferSize),
		UserID: userID,
	}

	re.Hub.Register <- client


	go client.WritePump()
	go client.ReadPump(re.handleWebSocketMessage)
}


func (re *Repository) handleWebSocketMessage(client *pkgwebsocket.Client, messageType int, data []byte) {
	if messageType != websocket.TextMessage {
		return
	}

	var msg pkgwebsocket.WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("failed to unmarshal WebSocket message: %v", err)
		return
	}

	switch msg.Type {
	case pkgwebsocket.MsgTypePrivateMsg:
		re.handlePrivateMessage(client, msg)
	default:
		log.Printf("unknown message type from user %s: %s", client.UserID, msg.Type)
	}
}

func (re *Repository) handlePrivateMessage(sender *pkgwebsocket.Client, msg pkgwebsocket.WSMessage) {

	var payload pkgwebsocket.PrivateMsgPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		log.Printf("failed to parse private_msg payload: %v", err)
		return
	}

	if payload.RecipientId == "" || payload.Text == "" {
		log.Printf("invalid private_msg payload: missing recipientId or text")
		return
	}


	savedMsg, err := re.MessageService.SendMessage(sender.UserID, payload.RecipientId, payload.Text)
	if err != nil {
		log.Printf("failed to save message: %v", err)
		return
	}


	senderNickname, err := re.MessageService.GetUserNickname(sender.UserID)
	if err != nil {
		senderNickname = sender.UserID
	}


	incomingPayload := pkgwebsocket.IncomingMsgPayload{
		MessageId:      savedMsg.MessageId,
		SenderId:       sender.UserID,
		SenderNickname: senderNickname,
		Text:           savedMsg.TextMessage,
		TimeStamp:      savedMsg.TimeStamp,
	}

	incomingData, err := json.Marshal(map[string]interface{}{
		"type":    pkgwebsocket.MsgTypeIncomingMsg,
		"payload": incomingPayload,
	})
	if err != nil {
		log.Printf("failed to marshal incoming message: %v", err)
		return
	}


	re.Hub.SendToUser(payload.RecipientId, incomingData)


	re.Hub.SendToUser(sender.UserID, incomingData)
}

func (re *Repository) SetHub(hub *pkgwebsocket.Hub) {
	re.Hub = hub
}