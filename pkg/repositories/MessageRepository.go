package repositories

import (
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
)

type MessageRepository interface {
	GetChatUsers(currentUserID string) ([]models.ChatUser, error)
}

func (db *DB) GetChatUsers(currentUserID string) ([]models.ChatUser, error) {
	query := `
		SELECT
			u.userId,
			u.nickName,
			(
				SELECT MAX(m.timeStamp)
				FROM message m
				WHERE (m.senderId = u.userId AND m.recipientId = ?) OR (m.senderId = ? AND m.recipientId = u.userId)
			) AS lastMessageTime
		FROM user u
		WHERE u.userId != ?
		ORDER BY
			CASE WHEN lastMessageTime IS NULL THEN 1 ELSE 0 END,
			lastMessageTime DESC,
			u.nickName ASC
	`

	rows, err := db.Conn.Query(query, currentUserID, currentUserID, currentUserID)
	if err != nil {
		return nil, realtimeforum.ErrInternal
	}
	defer rows.Close()

	var users []models.ChatUser
	for rows.Next() {
		var user models.ChatUser
		if err := rows.Scan(&user.UserId, &user.Nickname, &user.LastMessageTime); err != nil {
			return nil, realtimeforum.ErrInternal
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, realtimeforum.ErrInternal
	}

	if users == nil {
		users = []models.ChatUser{}
	}

	return users, nil
}