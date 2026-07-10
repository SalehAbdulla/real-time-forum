## Frontend SPA Implementation Plan

### 1. `static/` Directory Structure (to be created)
```
static/
├── css/
│   └── style.css
├── js/
│   ├── app.js           - Entry point, auth check, page initialization
│   ├── api.js           - Fetch wrapper with cookie credentials
│   ├── router.js        - Hash-based SPA router
│   ├── websocket.js     - WebSocket client
│   ├── auth.js          - Login/Register UI
│   ├── feed.js          - Post feed with pagination/sort
│   ├── postDetail.js    - Single post view + comments + reactions
│   ├── createPost.js    - Create post form
│   ├── chat.js          - Chat users list + conversation view
│   ├── notifications.js - Notifications list + mark-as-read
│   └── utils.js         - DOM helpers, date formatting, etc.
```

### 2. `templates/index.html` Update
Add `<link>` for CSS and `<script>` tags for all JS files (module scripts).

### 4. Frontend SPA Features

**Auth (cookie-based, already set by server):**
- `#/login` — Login form (identifier/password → form-encoded POST)
- `#/register` — Registration form (nickName, email, firstName, lastName, password, confirmPassword, age, gender → form-encoded POST)
- On successful auth, server sets `session_token` cookie → redirect to `#/feed`
- `#/logout` — POST /api/v1/auth/logout → redirect to `#/login`

**Router (hash-based):**
- `#/login` | `#/register` — Unauthenticated views
- `#/feed` — Post feed (default after login)
- `#/post/:id` — Post detail with comments & reactions
- `#/create` — Create post form
- `#/chat` — Chat user list
- `#/chat/:userId` — Chat with specific user
- `#/notifications` — Notifications page
- `#/profile` — User profile (from `/api/v1/auth/me`)
- Auth guard: check `GET /api/v1/auth/me` on init; redirect to login if 401

**Post Feed (`#/feed`):**
- Paginated list of posts (query params: page, size, sortBy, sortOrder)
- Sort by: Created At (default), Title, Score
- Each post shows: title, author, category, score, comment count, timestamp
- "Load More" button or infinite scroll
- Click post → navigate to `#/post/:id`

**Post Detail (`#/post/:id`):**
- Full post content
- Comments section with pagination
- Comment form (textarea → form-encoded POST)
- Reaction buttons (upvote/downvote) → `POST /api/v1/reactions`
- Real-time comment/reaction updates via WebSocket (if applicable)

**Create Post (`#/create`):**
- Form: title, content (textarea), category dropdown (fetched from `/api/v1/categories`)
- POST JSON body to `/api/v1/posts`
- On success → redirect to `#/feed`

**Chat (`#/chat`):**
- Left sidebar: list of users with online status (`GET /api/v1/messages/users`)
- Right panel: conversation with selected user (`GET /api/v1/messages?partnerId=X`)
- Message input → send via WebSocket (`private_msg`)
- Incoming messages → received via WebSocket (`incoming_msg`)
- User online/offline status → received via WebSocket (`user_status`)

**Notifications (`#/notifications`):**
- List of notifications with pagination
- Unread count badge in nav
- Mark single as read / mark all as read
- Real-time notifications via WebSocket (`notification`)

### 5. WebSocket Client
- Connect to `ws://localhost:3000/ws` (cookie-based auth)
- Handle: `incoming_msg`, `user_status`, `notification`
- Send: `private_msg` (with `{recipientId, text}` payload)
- Auto-reconnect on disconnect

### 6. CSS
- Clean, modern dark/Light theme
- Responsive layout: sidebar nav + main content area
- Mobile-friendly
