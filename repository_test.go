package main

import (
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestCreateOrUpdatePrompt_Unit(t *testing.T) {
	repo := NewFakePromptRepository()
	service := NewPromptService(repo)

	tests := []struct {
		name               string // description of this test case
		prompt             *Prompt
		want               *Prompt
		wantErr            bool
		failCreateOrUpdate bool
		failDelete         bool
		failGetByID        bool
		failGetAll         bool
	}{
		{
			name: "Create Prompt Success",
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
			wantErr:            false,
			failCreateOrUpdate: false,
		},
		{
			name: "Create Prompt Failure",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr:            true,
			failCreateOrUpdate: true,
		},

		{
			name: "Create Prompt Failure - Title is empty",
			prompt: &Prompt{
				Title:         "",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr:            true,
			failCreateOrUpdate: false,
		},
		{
			name: "Create Prompt Failure - Prompt Content is empty",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "",
			},
			wantErr:            true,
			failCreateOrUpdate: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.failCreateOrUpdate = tt.failCreateOrUpdate
			got, gotErr := service.CreateOrUpdatePrompt(tt.prompt)
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

type fakePromptRepository struct {
	prompts            map[int]*Prompt
	nextID             int
	failCreateOrUpdate bool
	failDelete         bool
	failGetByID        bool
	failGetAll         bool
}

// fake prompt repository
func NewFakePromptRepository() *fakePromptRepository {
	return &fakePromptRepository{
		prompts: make(map[int]*Prompt),
		nextID:  0,
	}
}

func (repo *fakePromptRepository) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	if prompt.Title == "" {
		return nil, errors.New("title is required")
	}

	if prompt.PromptContent == "" {
		return nil, errors.New("prompt content is required")
	}

	if repo.failCreateOrUpdate {
		return nil, errors.New("failed to create or update prompt")
	}
	repo.prompts[prompt.ID] = prompt
	repo.nextID++
	return prompt, nil
}

func (repo *fakePromptRepository) DeletePrompt(id int) error {
	if repo.failDelete {
		return errors.New("failed to delete prompt")
	}
	delete(repo.prompts, id)
	return nil
}

func (repo *fakePromptRepository) GetPromptByID(id int) (*Prompt, error) {
	if repo.failGetByID {
		return nil, errors.New("prompt not found")
	}
	return repo.prompts[id], nil
}

func (repo *fakePromptRepository) GetAllPrompts() ([]Prompt, error) {
	if repo.failGetAll {
		return nil, errors.New("failed to get all prompts")
	}
	var prompts []Prompt
	for _, prompt := range repo.prompts {
		prompts = append(prompts, *prompt)
	}
	return prompts, nil
}
