package repositories

import (
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
)

type CommentRepository interface {
	GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Comment, int, error)
}

func (db *DB) GetComments(postId int, pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Comment, int, error) {
	
	if err := db.DoesPostExists(postId); err != nil {
		return nil, 0, err
	}
	
	
	validSortColumns := map[string]string{
		"createdat": "c.createdAt",
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
		SELECT c.commentId, c.postId, c.userId, c.commentText, c.createdAt
		FROM comment c
		WHERE c.postId = ?
		ORDER BY ` + column + ` ` + order + `
		LIMIT ? OFFSET ?
	`

	rows, err := db.Conn.Query(query, postId, pageSize, offset)
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
			&com.CommentText,
			&com.CreatedAt,
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