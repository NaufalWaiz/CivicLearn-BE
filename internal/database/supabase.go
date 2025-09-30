package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("❌ DATABASE_URL is empty. Please set it in your .env")
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect DB: %v", err))
	}

	fmt.Println("✅ Connected to Supabase Postgres")
	DB = conn
}
