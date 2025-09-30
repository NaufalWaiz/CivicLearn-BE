package services

import (
	"github.com/google/uuid"
	"github.com/naufalwaiz/hukum-chatbot/config"
	"github.com/naufalwaiz/hukum-chatbot/internal/models"
)

// Create forum
func CreateForum(title, body, author string) (*models.Forum, error) {
	id := uuid.New().String()
	var forum models.Forum
	err := config.DB.QueryRow(
		"INSERT INTO forum (id, title, body, author) VALUES ($1,$2,$3,$4) RETURNING id, title, body, author, created_at, updated_at",
		id, title, body, author,
	).Scan(&forum.ID, &forum.Title, &forum.Body, &forum.Author, &forum.CreatedAt, &forum.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &forum, nil
}

// Get all forums
func GetAllForums() ([]models.Forum, error) {
	rows, err := config.DB.Query("SELECT id, title, body, author, created_at, updated_at FROM forum ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var f models.Forum
		if err := rows.Scan(&f.ID, &f.Title, &f.Body, &f.Author, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		forums = append(forums, f)
	}
	return forums, nil
}

// Get single forum by ID
func GetForumByID(id string) (*models.Forum, error) {
	var f models.Forum
	err := config.DB.QueryRow("SELECT id, title, body, author, created_at, updated_at FROM forum WHERE id=$1", id).
		Scan(&f.ID, &f.Title, &f.Body, &f.Author, &f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &f, nil
}

// Update forum
func UpdateForum(id, title, body, author string) error {
	_, err := config.DB.Exec(
		"UPDATE forum SET title=$1, body=$2, author=$3, updated_at=NOW() WHERE id=$4",
		title, body, author, id,
	)
	return err
}

// Delete forum
func DeleteForum(id string) error {
	_, err := config.DB.Exec("DELETE FROM forum WHERE id=$1", id)
	return err
}
