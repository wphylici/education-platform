package models

import (
	"time"
)

type Article struct {
	ID        int    `json:"id" gorm:"primaryKey;uniqueIndex"`
	Name      string `json:"name" gorm:"type:varchar(255);not null"`
	ChapterID int    `json:"chapterID" gorm:"type:integer;not null"`
	Content   string `json:"-" gorm:"type:text;not null"`

	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Chapter Chapter `json:"-"`
}

type CreateArticle struct {
	Name string `json:"name" binding:"required"`
}

type UpdateArticle struct {
	Name string `json:"name"`
}

type UpdateContent struct {
	Content string `json:"content"`
}

type ArticleResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ChapterID int    `json:"chapterID"`
	Route     string `json:"route"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
