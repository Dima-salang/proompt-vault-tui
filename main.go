package main

import (
	"log/slog"
	"os"

	"github.com/boltdb/bolt"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// repo := NewPromptRepository(db, logger)
	// service := NewPromptService(repo)
}

func openDB() (*bolt.DB, error) {
	return bolt.Open("prompts.db", 0600, nil)
}
