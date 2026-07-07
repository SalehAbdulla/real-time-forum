package repositories

import (
	"real-time-forum/pkg/models"
)

type PostRepository interface {
	GetPosts(pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Post, int, error)
}

func (db *DB) GetPosts(pageNumber int, pageSize int, sortBy string, sortOrder string) ([]models.Post, int, error) {
	validSortColumns := map[string]string{
		"createdat": "p.createdAt",
		"title":     "p.title",
		"score":     "p.score",
	}

	column, ok := validSortColumns[sortBy]
	if !ok {
		column = "p.createdAt"
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
	countQuery := "SELECT COUNT(*) FROM post"
	err := db.Conn.QueryRow(countQuery).Scan(&totalElements)
	if err != nil {
		return nil, 0, err
	}

	offset := (pageNumber - 1) * pageSize

	query := `
		SELECT p.postId, p.userId, u.nickName, p.title, p.content, 
			   p.categoryId, c.categoryName, p.score, p.commentsCounter, 
			   p.createdAt, p.updatedAt
		FROM post p
		JOIN user u ON p.userId = u.userId
		JOIN category c ON p.categoryId = c.categoryId
		ORDER BY ` + column + ` ` + order + `
		LIMIT ? OFFSET ?
	`

	rows, err := db.Conn.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostId,
			&post.UserId,
			&post.Nickname,
			&post.Title,
			&post.Content,
			&post.CategoryId,
			&post.CategoryName,
			&post.Score,
			&post.CommentsCounter,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	if posts == nil {
		posts = []models.Post{}
	}

	return posts, totalElements, nil
}