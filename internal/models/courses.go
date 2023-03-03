package models

import (
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type Course struct {
	ID          int       `gorm:"primaryKey;uniqueIndex"`
	Name        string    `gorm:"type:varchar(255);not null"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatorName string    `gorm:"type:varchar(255);not null"`
	ImageUrl    string    `gorm:"type:varchar(255)"`
	Category    string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	//Chapters    Chapters

	CreatedAt time.Time
	UpdatedAt time.Time

	Lecturer Lecturer `gorm:"foreignKey:CreatorID;references:UserID"`
}

type CreateCourse struct {
	Name        string         `form:"name" binding:"required"`
	CreatorID   string         `form:"creatorID" binding:"required,uuid"`
	CreatorName string         `form:"creatorName" binding:"required"`
	Image       multipart.File `form:"file,omitempty" validate:"required"`
	Category    string         `form:"category" binding:"required"`
	Description string         `form:"description" binding:"required"`
}

type CourseResponse struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	CreatorID   uuid.UUID `json:"creatorID" binding:"required"`
	CreatorName string    `json:"creatorName" binding:"required"`
	ImageUrl    string    `json:"imageUrl" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Route       string    `json:"route" binding:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Chapters struct {
}
