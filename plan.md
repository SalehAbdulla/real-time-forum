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

