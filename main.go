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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := vault.NewPromptRepository(db, logger)
	service := vault.NewPromptService(repo)

	p := tea.NewProgram(tui.NewModel(service), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("failed to run tui", "error", err)
		os.Exit(1)
	}
}

func openDB() (*bolt.DB, error) {
	return bolt.Open("prompts.db", 0600, nil)
}
