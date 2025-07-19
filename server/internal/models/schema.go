package models

import (
	"time"
)

var db = Db()
var DB = db

type Article struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AuthorID  string    `gorm:"column:authorID;not null;index" json:"authorID"`
	Title     string    `gorm:"column:title;not null;index" json:"title"`
	ImageUrl  string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
	ImagePath string    `gorm:"column:imagePath;default:null" json:"imagePath"`
	PostedAt  time.Time `gorm:"column:postedAt;index" json:"postedAt"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Author    *Author   `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"author,omitempty"`
}

type Author struct {
	ID         string     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name       string     `gorm:"column:name;unique;not null;index" json:"name"`
	AvatarUrl  string     `gorm:"column:avatarUrl;not null" json:"avatarUrl"`
	AvatarPath string     `gorm:"column:avatarPath;default:null" json:"avatarPath"`
	CreatedAt  time.Time  `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt;index" json:"updatedAt"`
	Article    []*Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"articles,omitempty"`
}
