Here is your updated and fully aligned blueprint. This revision incorporates the strict criteria of your project specification—specifically, **Single Page Application (SPA)** architecture via JavaScript, specialized **Private Messaging sort/pagination rules**, and our verified **SQLite PolyMorphic schema**.

---

# REST API Design — Real-Time Forum (SPA Aligned)

## General Conventions

| Convention | Standard |
| --- | --- |
| **Base URL** | `http://localhost:3000/api` |
| **Versioning** | `/api/v1/` |
| **Format** | JSON request/response bodies |
| **Auth** | Session-based cookie (or token handshake) |
| **Error Format** | `{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }` |
| **Success Format** | `{ "data": { ... } }` or `{ "data": [...] }` for lists |
| **Pagination** | `?offset=0&limit=10` for posts/comments; **Cursor/Page limits (`limit=10`)** for chat history scrolling. |

---

## 1. Authentication (Registration & Login)

*All required user profile data points are explicitly captured here.*

| # | API Name | Endpoint | Method | Purpose | Request Body | Response |
| --- | --- | --- | --- | --- | --- | --- |
| 1.1 | **Register** | `/api/v1/auth/register` | `POST` | Create a new user account | `{ "nickname": "string", "firstName": "string", "lastName": "string", "yearOfBirth": 1998, "gender": "string", "email": "string", "password": "string" }` | `201 Created` `{ "data": { "userId": 1, "nickname": "string" } }` (Sets session cookie) |
| 1.2 | **Login** | `/api/v1/auth/login` | `POST` | Authenticate user (nickname OR email + password) | `{ "login": "string", "password": "string" }` | `200 OK` `{ "data": { "userId": 1, "nickname": "string" } }` (Sets session cookie) |
| 1.3 | **Logout** | `/api/v1/auth/logout` | `POST` | Destroy session & update online status to `0` | — | `200 OK` `{ "data": { "message": "Logged out" } }` |
| 1.4 | **Check Session** | `/api/v1/auth/me` | `GET` | Used by SPA router on load to check login state | — | `200 OK` `{ "data": { "userId": 1, ... } }` or `401 Unauthorized` |

---

## 2. Categories & Posts

| # | API Name | Endpoint | Method | Purpose | Request Body | Response |
| --- | --- | --- | --- | --- | --- | --- |
| 2.1 | **List Categories** | `/api/v1/categories` | `GET` | Get categories for post creation forms | — | `200 OK` `{ "data": [ { "categoryId": 1, "categoryName": "Tech" } ] }` |
| 2.2 | **List Posts** | `/api/v1/posts` | `GET` | Fetch feed layout | — | `200 OK` `{ "data": [ { "postId": 1, "title": "...", "score": 0, "commentsCounter": 2 } ] }` |
| 2.3 | **Create Post** | `/api/v1/posts` | `POST` | Share new post under category | `{ "title": "string", "content": "string", "categoryId": 1 }` | `201 Created` |

---

## 3. Comments & Polymorphic Reactions

| # | API Name | Endpoint | Method | Purpose | Request Body | Response |
| --- | --- | --- | --- | --- | --- | --- |
| 3.1 | **List Comments** | `/api/v1/posts/{postId}/comments` | `GET` | Display comments dynamically when clicked | — | `200 OK` `{ "data": [ { "commentId": 1, "commentText": "..." } ] }` |
| 3.2 | **Create Comment** | `/api/v1/posts/{postId}/comments` | `POST` | Comment on a targeted post | `{ "commentText": "string" }` | `201 Created` |
| 3.3 | **React** | `/api/v1/reactions` | `POST` | Upsert a reaction vote via polymorphic parameters | `{ "entityType": "post/comment", "entityId": 1, "score": 1 }` | `200 OK` (Returns newly updated dynamic total score) |

---

## 4. Private Messages (Project-Specific Rules)

### Sorting Rules Applied:

1. Ordered by **last message time** (most recent contact first).
2. For new users with **no messages**, sorted **alphabetically by nickname**.

| # | API Name | Endpoint | Method | Purpose | Request Body / Params | Response |
| --- | --- | --- | --- | --- | --- | --- |
| 4.1 | **List Chat Users** | `/api/v1/messages/users` | `GET` | Returns list of all forum users with online/offline status, sorted strictly by specs. Always visible. | — | `200 OK` `{ "data": [ { "userId": 2, "nickname": "Alice", "isOnline": 1, "lastMessageTime": "..." } ] }` |
| 4.2 | **Get Chat Messages** | `/api/v1/messages/{userId}` | `GET` | **Loads last 10 messages** initially; client requests more using `offset` parameter when scrolling up. | `?offset=0&limit=10` | `200 OK` `{ "data": [ { "messageId": 45, "senderId": 1, "textMessage": "...", "timeStamp": "..." } ] }` |

---

## 5. WebSockets (Real-Time Communication)

| # | API Event Type | Direction | Purpose / Payload Structure |
| --- | --- | --- | --- |
| 5.1 | **Connection Up** | `Client → Server` | Open socket connection using handshake upgrade. |
| 5.2 | **Send Private Message** | `Client → Server` | `{ "type": "private_msg", "recipientId": 2, "text": "Hello" }` |
| 5.3 | **Broadcast Status** | `Server → Client` | `{ "type": "user_status", "userId": 5, "isOnline": 1 }` *(Fires immediately when users connect/disconnect)* |
| 5.4 | **Incoming Message** | `Server → Client` | `{ "type": "incoming_msg", "senderId": 2, "senderNickname": "Bob", "text": "Hello", "timeStamp": "..." }` |
# 6. Notifications

**Rules Applied** (matching your schema decisions):

- `entityType` is polymorphic — `comment` or `message` — resolved server-side to fetch the underlying post/comment/message details.
- Delivered two ways: WebSocket push if the recipient has a live connection (via the same hub as private messages), REST fetch on load/reconnect otherwise.
- Sorted by `createdAt DESC`, same cursor/offset pagination style as chat history.

| # | API Name | Endpoint | Method | Purpose | Request Body / Params | Response |
|---|----------|----------|--------|---------|------------------------|----------|
| 6.1 | List Notifications | `/api/v1/notifications` | GET | Fetch notification feed, newest first | `?offset=0&limit=10&unread=true` | `200 OK { "data": [ { "notificationId": 1, "actorId": "u2", "actorNickname": "Bob", "entityType": "comment", "entityId": 5, "isRead": 0, "createdAt": "..." } ] }` |
| 6.2 | Unread Count | `/api/v1/notifications/unread-count` | GET | Badge count for nav icon, polled or refreshed on WS event | — | `200 OK { "data": { "count": 4 } }` |
| 6.3 | Mark One Read | `/api/v1/notifications/{notificationId}/read` | PATCH | Mark a single notification as read (on click) | — | `200 OK { "data": { "notificationId": 1, "isRead": 1 } }` |
| 6.4 | Mark All Read | `/api/v1/notifications/read-all` | PATCH | Clear the badge in one action | — | `200 OK { "data": { "updated": 4 } }` |

# 5. WebSockets — add one event

| # | API Event Type | Direction | Purpose / Payload Structure |
|---|-----------------|-----------|------------------------------|
| 5.5 | New Notification | Server → Client | `{ "type": "notification", "notificationId": 9, "actorId": 2, "actorNickname": "Bob", "entityType": "comment", "entityId": 5, "createdAt": "..." }` — Fires when a comment/message is created for an online recipient — mirrors 5.4's shape. |





-----------------
Objectives
On this project you will have to focus on a few points:

Registration and Login

Creation of posts

Commenting posts

Private Messages

As you already did the first forum you can use part of the code, but not all of it. Your new forum will have five different parts:

SQLite, in which you will store data, just like in the previous forum

Golang, in which you will handle data and Websockets (Backend)

Javascript, in which you will handle all the Frontend events and clients Websockets

HTML, in which you will organize the elements of the page

CSS, in which you will stylize the elements of the page

You will have only one HTML file, so every change of page you want to do, should be handled in the Javascript. This can be called having a single page application.

Registration and Login
To be able to use the new and upgraded forum users will have to register and login, otherwise they will only see the registration or login page. This is premium stuff. The registration and login process should take in consideration the following features:

Users must be able to fill a register form to register into the forum. They will have to provide at least:

Nickname

Age

Gender

First Name

Last Name

E-mail

Password

The user must be able to connect using either the nickname or the e-mail combined with the password.

The user must be able to log out from any page on the forum.

Posts and Comments
This part is pretty similar to the first forum. Users must be able to:

Create posts

Posts will have categories as in the first forum

Create comments on the posts

See posts in a feed display

See comments only if they click on a post

Private Messages
Users will be able to send private messages to each other, so you will need to create a chat, where it will exist :

A section to show who is online/offline and able to talk to:

This section must be organized by the last message sent (just like discord). If the user is new and does not present messages you must organize it in alphabetic order.

The user must be able to send private messages to the users who are online.

This section must be visible at all times.

A section that when clicked on the user that you want to send a message, reloads the past messages. Chats between users must:

Be visible, for this you will have to be able to see the previous messages that you had with the user

Reload the last 10 messages and when scrolled up to see more messages you must provide the user with 10 more, without spamming the scroll event. Do not forget what you learned!! (Throttle, Debounce)

Messages must have a specific format:

A date that shows when the message was sent

The user name, that identifies the user that sent the message

As it is expected, the messages should work in real time, in other words, if a user sends a message, the other user should receive the notification of the new message without refreshing the page. Again this is possible through the usage of WebSockets in backend and frontend.

Allowed Packages
All standard go packages are allowed.

Gorilla websocket

sqlite3

bcrypt

gofrs/uuid or google/uuid


# Current Status & Next Steps

## ✅ Completed (REST API Backend):
1. **Authentication** — Register, Login, Logout, Session check
2. **Categories & Posts** — List categories, List/Create posts
3. **Comments & Reactions** — List/Create comments, React (polymorphic)
4. **Private Messages (REST)** — List chat users, Get chat messages (both implemented)

## ❌ Not Yet Implemented:
5. **WebSockets (Section 5)** — The entire real-time layer is missing:
   - No WebSocket upgrade endpoint
   - No `gorilla/websocket` dependency in `go.mod`
   - No handler for sending/receiving private messages in real-time
   - No online status broadcasting
   - No notification push via WebSocket

6. **Notifications (Section 6)** — REST endpoints not implemented
7. **Frontend (SPA)** — JavaScript single-page application not built yet

---

The **next logical implementation** is **Section 5: WebSockets (Real-Time Communication)**. Here's why:

1. WebSockets are the **core real-time feature** that makes this a "real-time forum"
2. They are a **prerequisite** for the frontend — the SPA needs WebSocket connections to show live messages and online status
3. Notifications (Section 6) depend on WebSockets for push delivery

## Proposed Implementation Plan for WebSockets:

### Step 1: Add `gorilla/websocket` dependency
```bash
go get github.com/gorilla/websocket
```

### Step 2: Create `pkg/websocket/` package with:
- **Hub** — Central manager that tracks all connected clients and routes messages
- **Client** — Represents a single WebSocket connection (read/write goroutines)
- **WebSocket handler** — Upgrade HTTP → WS, register client with Hub

### Step 3: Create message types for WebSocket communication:
- `private_msg` (Client → Server) — Send a private message
- `incoming_msg` (Server → Client) — Deliver received message to recipient
- `user_status` (Server → Client) — Broadcast online/offline changes
- `notification` (Server → Client) — Push new notifications
- `typing` (optional) — Typing indicators

### Step 4: Wire into existing architecture:
- Initialize Hub in `main.go`
- Add WebSocket route in `router.go` (e.g., `GET /ws`)
- Hook into `MessageService` to persist messages sent via WebSocket
- Hook into `SessionManager` to broadcast online status on connect/disconnect

### Step 5: Add `SaveMessage` method to `MessageRepository` (needed by WebSocket handler)

Would you like me to proceed with implementing the WebSocket layer? If so, please **toggle to ACT MODE** and I'll start building it step by step.

Here's a more detailed breakdown of the implementation tasks:

- [ ] Add `gorilla/websocket` dependency
- [ ] Create `pkg/websocket/hub.go` — Client connection registry & broadcast
- [ ] Create `pkg/websocket/client.go` — Read/Write goroutines per connection
- [ ] Create `pkg/websocket/types.go` — Message envelope types
- [ ] Add `SaveMessage` to `MessageRepository` interface & implementation
- [ ] Create WebSocket handler in `handlers` package
- [ ] Wire Hub into `main.go` and `Repository.go`
- [ ] Add WebSocket route to `router.go`
- [ ] Test with a simple WebSocket client (e.g., browser console)