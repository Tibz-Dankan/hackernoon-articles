package models

import (
	"time"
)

var db = Db()

type Article struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AuthorID      string    `gorm:"column:authorID;not null;index" json:"authorID"`
	Tag           string    `gorm:"column:tag;not null;index" json:"tag"`
	Title         string    `gorm:"column:title;not null;index" json:"title"`
	Href          string    `gorm:"column:href;default:null" json:"href"`
	ImageUrl      string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
	ImageFilename string    `gorm:"column:imageFilename;default:null" json:"imageFilename"`
	PostedAt      time.Time `gorm:"column:postedAt;index" json:"postedAt"`
	ReadDuration  string    `gorm:"column:readDuration" json:"readDuration"`
	CreatedAt     time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Author        *Author   `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"author,omitempty"`
}

type Author struct {
	ID             string     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name           string     `gorm:"column:name;unique;not null;index" json:"name"`
	AvatarUrl      string     `gorm:"column:avatarUrl;not null" json:"avatarUrl"`
	AvatarFilename string     `gorm:"column:avatarFilename;default:null" json:"avatarFilename"`
	PageUrl        string     `gorm:"column:pageUrl;default:null" json:"pageUrl"`
	CreatedAt      time.Time  `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt;index" json:"updatedAt"`
	Article        []*Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"articles,omitempty"`
}
