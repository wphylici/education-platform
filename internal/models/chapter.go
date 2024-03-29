package models

import "time"

type Chapter struct {
	ID       int    `json:"id" gorm:"primaryKey;uniqueIndex"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	CourseID int    `json:"courseID" gorm:"type:integer;not null"`

	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Course   Course    `json:"-"`
	Articles []Article `json:"articles" gorm:"constraint:OnDelete:CASCADE;"`
}

type CreateChapter struct {
	Name string `json:"name" binding:"required"`
}

type UpdateChapter struct {
	Name string `json:"name"`
}

type ChapterResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CourseID int    `json:"courseID"`
	Route    string `json:"route"`

	Articles []Article `json:"articles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
