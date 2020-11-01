package core

type Posts struct {
	Posts []*Post
}

type Post struct {
	Id      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}
