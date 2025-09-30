package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("SUPABASE_DB")
	if dsn == "" {
		log.Fatal("❌ SUPABASE_DB not set")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}

	// test connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	fmt.Println("✅ Connected to Supabase Postgres!")
}
