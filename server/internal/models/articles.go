package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *Article) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (a *Article) Create(article Article) (Article, error) {
	result := db.Create(&article)

	if result.Error != nil {
		return article, result.Error
	}
	return article, nil
}

func (a *Article) FindOne(id string) (Article, error) {
	var article Article
	db.First(&article, "id = ?", id)

	return article, nil
}

func (a *Article) FindByTitle(title string) (Article, error) {
	var article Article
	if err := db.First(&article, "title = ?", title).Error; err != nil {
		return article, err
	}

	return article, nil
}

func (a *Article) FindAll(limit float64, cursor string) ([]Article, error) {
	var articles []Article
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastArticle Article
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastArticle).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastArticle.UpdatedAt)
	}
	query.Where("\"imageUrl\" != ?", "imageUrl").Find(&articles)

	return articles, nil
}

func (a *Article) FindAllByPostedAt(limit float64, cursor string) ([]Article, error) {
	var articles []Article
	query := db.Order("\"postedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastArticle Article
		if err := db.Select("\"postedAt\"").Where("id = ?", cursor).First(&lastArticle).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"postedAt\" < ?", lastArticle.PostedAt)
	}
	query.Where("\"imageUrl\" != ?", "imageUrl").Find(&articles)

	return articles, nil
}

// TO add fetch by date functionality

func (a *Article) Search(searchQuery, cursor string, limit int) ([]Article, int64, error) {
	var articles []Article
	query := db.Model(&Article{}).Order("\"postedAt\" DESC").Limit(limit)

	query = query.Where("title ILIKE ? OR description ILIKE ?",
		"%"+searchQuery+"%", "%"+searchQuery+"%")

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if cursor != "" {
		var lastArticle Article
		if err := db.Model(&Article{}).Select("\"postedAt\"").Where("id = ?",
			cursor).First(&lastArticle).Error; err != nil {
			return nil, 0, err
		}
		query = query.Where("\"postedAt\" < ?", lastArticle.PostedAt)
	}

	if err := query.Find(&articles).Error; err != nil {
		return articles, 0, err
	}

	return articles, count, nil
}

func (a *Article) Update() (Article, error) {
	db.Save(&a)

	news, err := a.FindOne(a.ID)
	if err != nil {
		return news, err
	}

	return news, nil
}

func (a *Article) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Article{}).Error; err != nil {
		return err
	}

	return nil
}
