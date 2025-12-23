package main

import (
	"log/slog"
	"os"

	"github.com/Dima-salang/proompt-vault-tui/internal/vault"
	"github.com/Dima-salang/proompt-vault-tui/tui"
	"github.com/boltdb/bolt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// logging
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	logger := slog.New(slog.NewTextHandler(f, nil))

	// open the db connection
	db, err := openDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// create the repository and service
	repo := vault.NewPromptRepository(db, logger)
	service := vault.NewPromptService(repo)

	// run the tui
	p := tea.NewProgram(tui.NewModel(service), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("failed to run tui", "error", err)
		os.Exit(1)
	}
}

func openDB() (*bolt.DB, error) {
	return bolt.Open("prompts.db", 0600, nil)
}
