package repositories

import (
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
)

type CommentRepository interface {
	GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string, userID string) ([]models.Comment, int, error)
	CreateComment(userId string, postId int, content string) (models.Comment, error)
}

func (db *DB) GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string, userID string) ([]models.Comment, int, error) {

	if err := db.DoesPostExists(postId); err != nil {
		return nil, 0, err
	}

	validSortColumns := map[string]string{
		"createdat": "c.createdAt",
		"score":     "c.score",
	}

	column, ok := validSortColumns[sortBy]
	if !ok {
		column = "c.createdAt"
	}

	validSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	order, ok := validSortOrders[sortOrder]
	if !ok {
		order = "DESC"
	}

	var totalElements int
	countQuery := "SELECT COUNT(*) FROM comment WHERE postId = ?"
	err := db.Conn.QueryRow(countQuery, postId).Scan(&totalElements)
	if err != nil {
		return nil, 0, realtimeforum.ErrInternal
	}

	offset := (pageNumber - 1) * pageSize

	query := `
		SELECT c.commentId, c.postId, c.userId, u.nickName, c.commentText, c.score, c.createdAt,
		       COALESCE(r.score, 0) AS userScore
		FROM comment c
		JOIN user u ON c.userId = u.userId
		LEFT JOIN reaction r ON r.entityType = 'comment' AND r.entityId = c.commentId AND r.userId = ?
		WHERE c.postId = ?
		ORDER BY ` + column + ` ` + order + `
		LIMIT ? OFFSET ?
	`

	rows, err := db.Conn.Query(query, userID, postId, pageSize, offset)
	if err != nil {
		return nil, 0, realtimeforum.ErrInternal
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var com models.Comment
		err := rows.Scan(
			&com.CommentId,
			&com.PostId,
			&com.UserId,
			&com.Nickname,
			&com.CommentText,
			&com.Score,
			&com.CreatedAt,
			&com.UserScore,
		)
		if err != nil {
			return nil, 0, realtimeforum.ErrInternal
		}
		comments = append(comments, com)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, realtimeforum.ErrInternal
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, totalElements, nil
}

func (db *DB) CreateComment(userId string, postId int, content string) (models.Comment, error) {
	if err := db.DoesPostExists(postId); err != nil {
		return models.Comment{}, err
	}

	result, err := db.Conn.Exec(
		`INSERT INTO comment (postId, userId, commentText, score, createdAt)
		 VALUES (?, ?, ?, 0, datetime('now'))`,
		postId, userId, content,
	)
	if err != nil {
		return models.Comment{}, realtimeforum.ErrInternal
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, realtimeforum.ErrInternal
	}

	// Update the comment count on the post
	_, err = db.Conn.Exec(
		`UPDATE post SET commentsCounter = commentsCounter + 1 WHERE postId = ?`,
		postId,
	)
	if err != nil {
		return models.Comment{}, realtimeforum.ErrInternal
	}

	var com models.Comment
	err = db.Conn.QueryRow(
		`SELECT commentId, postId, userId, commentText, createdAt
		 FROM comment
		 WHERE commentId = ?`, commentID,
	).Scan(
		&com.CommentId,
		&com.PostId,
		&com.UserId,
		&com.CommentText,
		&com.CreatedAt,
	)
	
	if err != nil {
		return models.Comment{}, realtimeforum.ErrInternal
	}

	return com, nil
}
