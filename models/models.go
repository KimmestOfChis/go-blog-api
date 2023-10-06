package models

type Post struct {
	ID int `json:"-"`
	Title string `json:"title"`
	Content string `json:"content"`
}

