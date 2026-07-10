package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"real-time-forum/pkg/app/handlers"
	"real-time-forum/pkg/app/repositories"
	"real-time-forum/pkg/app/service"
	"real-time-forum/pkg/config"
	"real-time-forum/pkg/logger"
	pkgmiddleware "real-time-forum/pkg/middleware"
	"real-time-forum/pkg/render"
	pkgwebsocket "real-time-forum/pkg/websocket"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// testSchema is the SQL schema used to set up the in-memory database for testing.
const testSchema = `
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS reaction;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS user;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

VACUUM;
PRAGMA foreign_keys = ON;

CREATE TABLE user (
    userId          TEXT PRIMARY KEY,
    nickName        TEXT NOT NULL UNIQUE,
    firstName       TEXT NOT NULL,
    lastName        TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    hashedPassword  TEXT NOT NULL,
    yearOfBirth     INTEGER NOT NULL,
    gender          TEXT CHECK(gender IN ('male', 'female')),
    createdAt       TEXT DEFAULT (CURRENT_TIMESTAMP),
    updatedAt       TEXT DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE category (
    categoryId      INTEGER PRIMARY KEY AUTOINCREMENT,
    categoryName    TEXT NOT NULL UNIQUE
);

insert into category (categoryName) values ("General");
insert into category (categoryName) values ("Tech");
insert into category (categoryName) values ("Dev");
insert into category (categoryName) values ("Gaming");
insert into category (categoryName) values ("Help");
insert into category (categoryName) values ("Life");
insert into category (categoryName) values ("Sport");
insert into category (categoryName) values ("Misc");

CREATE TABLE post (
    postId           INTEGER PRIMARY KEY AUTOINCREMENT,
    userId           TEXT NOT NULL,
    title            TEXT NOT NULL,
    content          TEXT NOT NULL,
    categoryId       INTEGER NOT NULL,
    score            INTEGER DEFAULT 0,
    commentsCounter  INTEGER DEFAULT 0,
    createdAt        TEXT DEFAULT (CURRENT_TIMESTAMP),
    updatedAt        TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (userId) REFERENCES user(userId) ON DELETE CASCADE,
    FOREIGN KEY (categoryId) REFERENCES category(categoryId) ON DELETE RESTRICT
);

CREATE TABLE comment (
    commentId   INTEGER PRIMARY KEY AUTOINCREMENT,
    postId      INTEGER NOT NULL,
    userId      TEXT NOT NULL,
    commentText TEXT NOT NULL,
    score       INTEGER DEFAULT 0,                
    createdAt   TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (postId) REFERENCES post(postId) ON DELETE CASCADE,
    FOREIGN KEY (userId) REFERENCES user(userId) ON DELETE CASCADE
);

CREATE TABLE reaction (
    reactionId  INTEGER PRIMARY KEY AUTOINCREMENT,
    userId      TEXT NOT NULL,
    entityType  TEXT NOT NULL CHECK(entityType IN ('post', 'comment')),
    entityId    INTEGER NOT NULL,
    score       INTEGER NOT NULL CHECK(score IN (1, -1)),
    createdAt   TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (userId) REFERENCES user(userId) ON DELETE CASCADE,
    
    UNIQUE(userId, entityType, entityId)
);

CREATE TABLE session (
    sessionToken TEXT PRIMARY KEY,
    userId       TEXT NOT NULL UNIQUE,
    timeStamp    TEXT DEFAULT (CURRENT_TIMESTAMP),
    createdAt    TEXT DEFAULT (CURRENT_TIMESTAMP),
    expiredAt    TEXT NOT NULL,
    FOREIGN KEY (userId) REFERENCES user(userId) ON DELETE CASCADE
);

CREATE TABLE message (
    messageId   INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId    TEXT NOT NULL,
    recipientId TEXT NOT NULL,
    textMessage TEXT NOT NULL,
    timeStamp   TEXT DEFAULT (CURRENT_TIMESTAMP),
    isRead      INTEGER DEFAULT 0 CHECK(isRead IN (0, 1)),
    FOREIGN KEY (senderId) REFERENCES user(userId) ON DELETE CASCADE,
    FOREIGN KEY (recipientId) REFERENCES user(userId) ON DELETE CASCADE
);

CREATE TABLE notification (
    notificationId  INTEGER PRIMARY KEY AUTOINCREMENT,
    userId          TEXT NOT NULL,
    actorId         TEXT,
    entityType      TEXT NOT NULL CHECK(entityType IN ('comment', 'message')),
    entityId        INTEGER NOT NULL,
    isRead          INTEGER DEFAULT 0 CHECK(isRead IN (0, 1)),
    createdAt       TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (userId) REFERENCES user(userId) ON DELETE CASCADE,
    FOREIGN KEY (actorId) REFERENCES user(userId) ON DELETE SET NULL
);

CREATE INDEX idx_messages_chat_flow 
ON message (senderId, recipientId, timeStamp DESC);

CREATE INDEX idx_reactions_lookup 
ON reaction (entityType, entityId);

CREATE INDEX idx_posts_category 
ON post (categoryId);

CREATE INDEX idx_post_userId ON post (userId);
CREATE INDEX idx_comment_postId ON comment (postId);
CREATE INDEX idx_comment_userId ON comment (userId);

CREATE INDEX idx_notifications_user 
ON notification (userId, isRead, createdAt DESC);
`

// TestWebSocketPrivateMessage tests the full WebSocket flow:
// 1. Register two users
// 2. Open WebSocket connections for both
// 3. Send a private message from user1 to user2
// 4. Verify the message was saved in the database via REST API
func TestWebSocketPrivateMessage(t *testing.T) {
	// ---------- Setup ----------
	origDir, _ := os.Getwd()
	os.Chdir("..")
	defer os.Chdir(origDir)

	app = config.AppConfig{}
	app.InProduction = false
	app.UseCache = false

	logger.InitLogger(&app)

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		t.Fatalf("failed to create template cache: %v", err)
	}
	app.TemplateCache = templateCache
	render.NewTemplates(&app)

	// Use in-memory database for testing
	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	defer database.Close()

	_, err = database.Exec(testSchema)
	if err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	// Initialize services and handlers
	dbConn := &repositories.DB{Conn: database}
	authService := service.NewAuthService(dbConn)
	messageService := service.NewMessageService(dbConn, dbConn)

	hc := handlers.NewHandlerContext(&app, authService, nil, nil, nil, nil, messageService)
	handlers.SetHandlerContext(hc)

	wsHub := pkgwebsocket.NewHub()
	hc.SetHub(wsHub)
	go wsHub.Run()

	// Build a handler that routes auth endpoints directly and protected endpoints through AuthMiddleware.
	wsHandler := pkgmiddleware.AuthMiddleware(http.HandlerFunc(handlers.HandlerCtx.ServeWs))
	messagesHandler := pkgmiddleware.AuthMiddleware(http.HandlerFunc(handlers.HandlerCtx.GetChatMessages))
	mainMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/register":
			handlers.HandlerCtx.Register(w, r)
		case "/api/v1/auth/login":
			handlers.HandlerCtx.Login(w, r)
		case "/api/v1/messages":
			messagesHandler.ServeHTTP(w, r)
		case "/ws":
			wsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	server := &http.Server{Handler: mainMux}
	go server.Serve(listener)
	defer server.Close()
	defer listener.Close()

	baseURL := "http://" + listener.Addr().String()
	wsURL := "ws://" + listener.Addr().String() + "/ws"

	// ---------- Register two users ----------
	user1ID, token1 := registerTestUser(t, baseURL, "testuser1", "user1@test.com", "TestPass123!@#")
	user2ID, token2 := registerTestUser(t, baseURL, "testuser2", "user2@test.com", "TestPass123!@#")

	t.Logf("user1: id=%s, token=%s", user1ID, token1)
	t.Logf("user2: id=%s, token=%s", user2ID, token2)

	// ---------- Open WebSocket connections ----------
	ws1 := dialWebSocket(t, wsURL, token1)
	defer ws1.Close()

	ws2 := dialWebSocket(t, wsURL, token2)
	defer ws2.Close()

	// Give the hub time to register both clients
	time.Sleep(200 * time.Millisecond)

	// ---------- user1 sends private message to user2 ----------
	privateMsg := map[string]interface{}{
		"type": "private_msg",
		"payload": map[string]string{
			"recipientId": user2ID,
			"text":        "Hello, user2!",
		},
	}
	msgBytes, _ := json.Marshal(privateMsg)
	err = ws1.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		t.Fatalf("ws1 write error: %v", err)
	}

	// Give the server time to process the message
	time.Sleep(200 * time.Millisecond)

	// ---------- Verify the message was saved via REST API ----------
	messagesResp := getMessages(t, baseURL, token1, user2ID, 0, 10)
	msgs := messagesResp.Data.Messages
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}

	msg := msgs[0]
	if msg.SenderId != user1ID {
		t.Fatalf("expected senderId %s, got %s", user1ID, msg.SenderId)
	}
	if msg.RecipientId != user2ID {
		t.Fatalf("expected recipientId %s, got %s", user2ID, msg.RecipientId)
	}
	if msg.TextMessage != "Hello, user2!" {
		t.Fatalf("expected text 'Hello, user2!', got '%s'", msg.TextMessage)
	}
	t.Logf("message verified: sender=%s, recipient=%s, text=%s", msg.SenderId, msg.RecipientId, msg.TextMessage)

	// ---------- Also verify from user2's perspective ----------
	messagesResp2 := getMessages(t, baseURL, token2, user1ID, 0, 10)
	msgs2 := messagesResp2.Data.Messages
	if len(msgs2) != 1 {
		t.Fatalf("user2: expected 1 message, got %d", len(msgs2))
	}
	t.Log("message visible from both user perspectives: PASS")
}

// messagesResponse mirrors the API response structure for messages.
type messagesResponse struct {
	Data struct {
		Messages      []messageDTO `json:"messages"`
		Offset        int          `json:"offset"`
		Limit         int          `json:"limit"`
		TotalElements int          `json:"totalElements"`
	} `json:"data"`
}

type messageDTO struct {
	MessageId   int    `json:"messageId"`
	SenderId    string `json:"senderId"`
	RecipientId string `json:"recipientId"`
	TextMessage string `json:"textMessage"`
	TimeStamp   string `json:"timeStamp"`
	IsRead      int    `json:"isRead"`
}

// getMessages fetches messages from the REST API.
func getMessages(t *testing.T, baseURL, token, partnerID string, offset, limit int) messagesResponse {
	t.Helper()
	url := fmt.Sprintf("%s/api/v1/messages?partnerId=%s&offset=%d&limit=%d", baseURL, partnerID, offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Cookie", "session_token="+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get messages request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		t.Fatalf("get messages expected 200, got %d: %s", resp.StatusCode, body.String())
	}

	var result messagesResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

// registerTestUser registers a user and returns (userID, sessionToken).
func registerTestUser(t *testing.T, baseURL, nickname, email, password string) (string, string) {
	t.Helper()

	form := fmt.Sprintf("nickName=%s&email=%s&firstName=Test&lastName=User&password=%s&confirmPassword=%s&age=25&gender=male",
		nickname, email, password, password)

	resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/x-www-form-urlencoded", bytes.NewBufferString(form))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		resp.Body.Close()
		t.Fatalf("register expected 201, got %d: %s", resp.StatusCode, body.String())
	}

	var token string
	for _, c := range resp.Cookies() {
		if c.Name == "session_token" {
			token = c.Value
			break
		}
	}
	if token == "" {
		resp.Body.Close()
		t.Fatal("no session_token cookie in register response")
	}

	var result struct {
		UserID  string `json:"userId"`
		Message string `json:"message"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	return result.UserID, token
}

// dialWebSocket opens a WebSocket connection with the session token as a cookie.
func dialWebSocket(t *testing.T, wsURL, token string) *websocket.Conn {
	t.Helper()

	header := http.Header{}
	header.Add("Cookie", "session_token="+token)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		t.Fatalf("websocket dial error: %v", err)
	}
	return conn
}