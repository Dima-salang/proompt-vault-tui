package main

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestCreateOrUpdatePrompt_Integration(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		prompt  *Prompt
		want    *Prompt
		wantErr bool
	}{
		{
			name: "test case 1",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			want: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temporary DB
			dir := t.TempDir()
			dbPath := filepath.Join(dir, "test.db")
			db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			repo := NewPromptRepository(db, logger)

			got, gotErr := repo.CreateOrUpdatePrompt(tt.prompt)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateOrUpdatePrompt() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateOrUpdatePrompt() succeeded unexpectedly")
			}

			// compare the prompts
			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("CreateOrUpdatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
