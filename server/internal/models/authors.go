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

func (a *Author) Create(author Author) (Author, error) {
	result := db.Create(&author)

	if result.Error != nil {
		return author, result.Error
	}
	return author, nil
}

func (a *Author) FindOne(id string) (Author, error) {
	var author Author
	db.First(&author, "id = ?", id)

	return author, nil
}

func (a *Author) FindByName(name string) (Author, error) {
	var author Author
	if err := db.First(&author, "name = ?", name).Error; err != nil {
		return author, err
	}

	return author, nil
}

func (a *Author) FindByPage(pageURL string) (Author, error) {
	var author Author
	if err := db.First(&author, "\"pageURL\" = ?", pageURL).Error; err != nil {
		return author, err
	}

	return author, nil
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

func (a *Author) Update() (Author, error) {
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
