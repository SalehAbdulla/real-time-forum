package models

type Comment struct {
	CommentId   int    `json:"commentId"`
	PostId      int    `json:"postId"`
	UserId      string `json:"userId"`
	Nickname    string `json:"nickname"`
	CommentText string `json:"commentText"`
	Score       int    `json:"score"`
	CreatedAt   string `json:"createdAt"`
}
