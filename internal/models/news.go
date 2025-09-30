package models

import "time"

type News struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Author    string     `json:"author"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
