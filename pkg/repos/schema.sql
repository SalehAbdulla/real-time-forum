-- Disable foreign key checks temporarily to safely drop tables if they exist
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

PRAGMA foreign_keys = ON;

-- 1. USERS TABLE
CREATE TABLE users (
    userId          INTEGER PRIMARY KEY AUTOINCREMENT,
    nickName        TEXT NOT NULL UNIQUE,
    firstName       TEXT NOT NULL,
    lastName        TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    hashedPassword  TEXT NOT NULL,
    yearOfBirth     INTEGER NOT NULL, -- Stored directly (Age is derived dynamically)
    gender          TEXT CHECK(gender IN ('male', 'female')),
    isOnline        INTEGER DEFAULT 0 CHECK(isOnline IN (0, 1)),
    createdAt       TEXT DEFAULT (CURRENT_TIMESTAMP),
    updatedAt       TEXT DEFAULT (CURRENT_TIMESTAMP)
);

-- 2. CATEGORIES TABLE
CREATE TABLE categories (
    categoryId      INTEGER PRIMARY KEY AUTOINCREMENT,
    categoryName    TEXT NOT NULL UNIQUE,
    createdAt       TEXT DEFAULT (CURRENT_TIMESTAMP),
    updatedAt       TEXT DEFAULT (CURRENT_TIMESTAMP)
);

-- 3. POSTS TABLE
CREATE TABLE posts (
    postId           INTEGER PRIMARY KEY AUTOINCREMENT,
    userId           INTEGER NOT NULL,
    title            TEXT NOT NULL,
    content          TEXT NOT NULL,
    categoryId       INTEGER NOT NULL,
    score            INTEGER DEFAULT 0,            -- Cached total score (upvotes - downvotes)
    commentsCounter  INTEGER DEFAULT 0,            -- Cached total number of comments
    createdAt        TEXT DEFAULT (CURRENT_TIMESTAMP),
    updatedAt        TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
    FOREIGN KEY (categoryId) REFERENCES categories(categoryId) ON DELETE RESTRICT
);

-- 4. COMMENTS TABLE
CREATE TABLE comments (
    commentId   INTEGER PRIMARY KEY AUTOINCREMENT,
    postId      INTEGER NOT NULL,
    userId      INTEGER NOT NULL,
    commentText TEXT NOT NULL,
    score       INTEGER DEFAULT 0,                 -- Cached score for the comment
    createdAt   TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE,
    FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE
);

-- 5. POLYMORPHIC REACTIONS TABLE (Abstracted for Posts & Comments)
CREATE TABLE reactions (
    reactionId  INTEGER PRIMARY KEY AUTOINCREMENT,
    userId      INTEGER NOT NULL,
    entityType  TEXT NOT NULL CHECK(entityType IN ('post', 'comment')),
    entityId    INTEGER NOT NULL,
    score       INTEGER NOT NULL CHECK(score IN (1, -1)), -- 1 = Like, -1 = Dislike
    createdAt   TEXT DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
    
    -- Crucial Constraint: Ensures a user can only vote on a specific item once
    UNIQUE(userId, entityType, entityId)
);

-- 6. SESSIONS TABLE
CREATE TABLE sessions (
    sessionToken TEXT PRIMARY KEY, -- Clean lookup token as the unique identifier
    userId       INTEGER NOT NULL,
    timeStamp    TEXT DEFAULT (CURRENT_TIMESTAMP),
    createdAt    TEXT DEFAULT (CURRENT_TIMESTAMP),
    expiredAt    TEXT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE
);

-- 7. MESSAGES TABLE (Direct Messaging Architecture)
CREATE TABLE messages (
    messageId   INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId    INTEGER NOT NULL,
    recipientId INTEGER NOT NULL,
    textMessage TEXT NOT NULL,
    timeStamp   TEXT DEFAULT (CURRENT_TIMESTAMP),
    isRead      INTEGER DEFAULT 0 CHECK(isRead IN (0, 1)),
    FOREIGN KEY (senderId) REFERENCES users(userId) ON DELETE CASCADE,
    FOREIGN KEY (recipientId) REFERENCES users(userId) ON DELETE CASCADE
);

-- PERFORMANCE & LOOKUP INDEXES

-- Index for real-time direct messaging chat timelines
CREATE INDEX idx_messages_chat_flow 
ON messages (senderId, recipientId, timeStamp DESC);

-- Indexes for polymorphic reaction calculations
CREATE INDEX idx_reactions_lookup 
ON reactions (entityType, entityId);

-- Index for retrieving target feed posts quickly
CREATE INDEX idx_posts_category 
ON posts (categoryId);