package models

import (
	"time"

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

func (a *Article) FindAll(limit float64, cursor string) ([]Article, int64, error) {
	var articles []Article
	var count int64
	query := db.Model(&Article{}).
		// Preload("Author").
		Order("\"updatedAt\" DESC").
		Limit(int(limit))

	if cursor != "" {
		var lastArticle Article
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastArticle).Error; err != nil {
			return nil, 0, err
		}
		query = query.Where("\"updatedAt\" < ?", lastArticle.UpdatedAt)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query.Find(&articles)

	return articles, count, nil
}

func (a *Article) FindAllByPostedAt(limit int, articleIDCursor string,
	dateCursor time.Time, offset int) ([]Article, int64, error) {
	var articles []Article
	var count int64
	query := db.Model(&Article{}).
		Preload("Author").
		Order("\"postedAt\" DESC").
		Limit(int(limit))

	if offset != 0 {
		query = query.Offset(offset)
	}

	if articleIDCursor != "" {
		var lastArticle Article
		if err := db.Select("\"postedAt\"").Where("id = ?",
			articleIDCursor).First(&lastArticle).Error; err != nil {
			return nil, 0, err
		}
		query = query.Where("\"postedAt\" < ?", lastArticle.PostedAt)
	}

	if !dateCursor.IsZero() {
		query = query.Where("\"postedAt\" <= ?", dateCursor)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query.Find(&articles)

	return articles, count, nil
}

func (a *Article) Search(searchQuery, articleIDCursor string,
	dateCursor time.Time, limit int, offset int) ([]Article, int64, error) {
	var articles []Article
	query := db.Model(&Article{}).
		Preload("Author").
		Order("\"postedAt\" DESC").
		Limit(limit)

	query = query.Where("title ILIKE ?", "%"+searchQuery+"%")

	if offset != 0 {
		query = query.Offset(offset)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if articleIDCursor != "" {
		var lastArticle Article
		if err := db.Model(&Article{}).Select("\"postedAt\"").Where("id = ?",
			articleIDCursor).First(&lastArticle).Error; err != nil {
			return nil, 0, err
		}
		query = query.Where("\"postedAt\" < ?", lastArticle.PostedAt)
	}

	if !dateCursor.IsZero() {
		query = query.Where("\"postedAt\" <= ?", dateCursor)
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
