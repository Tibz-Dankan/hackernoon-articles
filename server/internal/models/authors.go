package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *Author) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (a *Author) Create(article Article) (Article, error) {
	result := db.Create(&article)

	if result.Error != nil {
		return article, result.Error
	}
	return article, nil
}

func (a *Author) FindOne(id string) (Article, error) {
	var article Article
	db.First(&article, "id = ?", id)

	return article, nil
}

func (a *Author) FindByTitle(title string) (Article, error) {
	var article Article
	if err := db.First(&article, "title = ?", title).Error; err != nil {
		return article, err
	}

	return article, nil
}

func (a *Author) FindAll(limit float64, cursor string) ([]Author, error) {
	var authors []Author
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastAuthor Author
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastAuthor).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastAuthor.UpdatedAt)
	}
	query.Find(&authors)

	return authors, nil
}

func (a *Author) Update() (Article, error) {
	db.Save(&a)

	author, err := a.FindOne(a.ID)
	if err != nil {
		return author, err
	}

	return author, nil
}

func (a *Author) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Author{}).Error; err != nil {
		return err
	}

	return nil
}
