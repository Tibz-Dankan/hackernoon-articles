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

func (a *Article) FindAllByPostedAtInAsc(limit int) ([]Article, int64, error) {
	var articles []Article
	var count int64
	query := db.Model(&Article{}).
		Order("\"postedAt\" ASC").
		Limit(int(limit))

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	query.Find(&articles)

	return articles, count, nil
}

func (a *Article) FindAllWithWrongImage(limit float64, cursor string) ([]Article, int64, error) {
	var articles []Article
	query := db.Model(&Article{}).
		Order("\"postedAt\" DESC").
		Limit(int(limit))

	wrongFormat := "?auto"

	query = query.Where("\"imageUrl\" ILIKE ?", "%"+wrongFormat+"%")

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Find(&articles).Error; err != nil {
		return articles, 0, err
	}

	return articles, count, nil
}

func (a *Article) FindCount() (int64, error) {
	var count int64
	if err := db.Model(&Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
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

func (a *Article) SearchByTagIndex(searchQuery string) ([]Article, int64, error) {
	var articles []Article
	query := db.Model(&Article{}).
		Preload("Author").
		Order("\"postedAt\" DESC")

	// query = query.Where("\"tagIndex\" ILIKE ?", "%"+searchQuery+"%")
	query = query.Where("\"tagIndex\" = ?", searchQuery)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Find(&articles).Error; err != nil {
		return articles, 0, err
	}

	return articles, count, nil
}

func (a *Article) FindByPostedAt(date time.Time) ([]Article, error) {
	var articles []Article
	
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	
	err := db.Model(&Article{}).
		Preload("Author").
		Where("\"postedAt\" >= ? AND \"postedAt\" < ?", startOfDay, endOfDay).
		Order("\"postedAt\" DESC").
		Find(&articles).Error
	
	if err != nil {
		return nil, err
	}
	
	return articles, nil
}

func (a *Article) FindArticleCountPerDay(limit int, dateCursor time.Time) ([]map[string]interface{}, error) {
	var results []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}

	query := db.Model(&Article{}).
		Select("DATE(\"postedAt\") as date, COUNT(*) as count").
		Group("DATE(\"postedAt\")").
		Order("date DESC").
		Limit(limit)

	if !dateCursor.IsZero() {
		query = query.Where("DATE(\"postedAt\") < DATE(?)", dateCursor)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Create a map for quick lookup of existing dates and their counts
	dateCountMap := make(map[string]int64)
	for _, result := range results {
		dateStr := result.Date.Format("2006-01-02")
		dateCountMap[dateStr] = result.Count
	}

	// Get the first (latest) and last (oldest) dates from results
	firstDate := results[0].Date             // Latest date
	lastDate := results[len(results)-1].Date // Oldest date

	// Fill in missing dates between first and last date
	var dayArticleCounts []map[string]interface{}
	currentDate := firstDate

	for !currentDate.Before(lastDate) {
		dateStr := currentDate.Format("2006-01-02")

		// Check if date exists in our results
		count, exists := dateCountMap[dateStr]
		if !exists {
			count = 0 // Set count to 0 for missing dates
		}

		dayData := map[string]interface{}{
			"date":  dateStr,
			"count": count,
		}
		dayArticleCounts = append(dayArticleCounts, dayData)

		// Move to previous day (since we're in descending order)
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	return dayArticleCounts, nil
}


func (a *Article) CountDistinctDays(count *int64) error {
	return db.Model(&Article{}).
		Select("COUNT(DISTINCT DATE(\"postedAt\"))").
		Row().Scan(count)
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
