package models

type Comment struct {
	CommentId   int    `json:"commentId"`
	PostId      int    `json:"postId"`
	UserId      string `json:"userId"`
	CommentText string `json:"commentText"`
	CreatedAt   string `json:"createdAt"`
}