package services

import (
	"github.com/google/uuid"
	"github.com/naufalwaiz/hukum-chatbot/config"
	"github.com/naufalwaiz/hukum-chatbot/internal/models"
)

// Create news
func CreateNews(title, content, author string) (string, string, error) {
	var id, returnedAuthor string
	err := config.DB.QueryRow(
		"INSERT INTO news (id, title, content, author) VALUES ($1,$2,$3,$4) RETURNING id, author",
		uuid.New().String(), title, content, author,
	).Scan(&id, &returnedAuthor)
	if err != nil {
		return "", "", err
	}
	return id, returnedAuthor, nil
}

// Get all news
func GetAllNews() ([]models.News, error) {
	rows, err := config.DB.Query(`
		SELECT id, title, content, author, created_at, updated_at 
		FROM news ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Author, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

// Get single news
func GetNewsByID(id string) (*models.News, error) {
	var n models.News
	err := config.DB.QueryRow(`
		SELECT id, title, content, author, created_at, updated_at 
		FROM news WHERE id=$1
	`, id).Scan(&n.ID, &n.Title, &n.Content, &n.Author, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &n, nil
}

// Update news
func UpdateNews(id, title, content, author string) error {
	_, err := config.DB.Exec(
		"UPDATE news SET title=$1, content=$2, author=$3, updated_at=NOW() WHERE id=$4",
		title, content, author, id,
	)
	return err
}

// Delete news
func DeleteNews(id string) error {
	_, err := config.DB.Exec("DELETE FROM news WHERE id=$1", id)
	return err
}
