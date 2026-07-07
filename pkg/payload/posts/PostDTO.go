package posts


type Score int

const (
	Neutral Score = iota
	Like
	Dislike
)

type PostDTO struct {
	postId int
	userId string
	title string
	content string 
	categoryId int
	createdAt string
	updatedAt string
	score Score
}