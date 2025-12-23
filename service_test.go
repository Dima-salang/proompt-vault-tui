package main

import (
	"errors"
	"testing"
)

type MockPromptRepository struct {
	createOrUpdatePromptFunc func(prompt *Prompt) (*Prompt, error)
	deletePromptFunc         func(id int) error
	getPromptByIDFunc        func(id int) (*Prompt, error)
	getAllPromptsFunc        func() ([]Prompt, error)
}

func (m *MockPromptRepository) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	if m.createOrUpdatePromptFunc != nil {
		return m.createOrUpdatePromptFunc(prompt)
	}
	return prompt, nil
}

func (m *MockPromptRepository) DeletePrompt(id int) error {
	if m.deletePromptFunc != nil {
		return m.deletePromptFunc(id)
	}
	return nil
}

func (m *MockPromptRepository) GetPromptByID(id int) (*Prompt, error) {
	if m.getPromptByIDFunc != nil {
		return m.getPromptByIDFunc(id)
	}
	return nil, nil
}

func (m *MockPromptRepository) GetAllPrompts() ([]Prompt, error) {
	if m.getAllPromptsFunc != nil {
		return m.getAllPromptsFunc()
	}
	return nil, nil
}

func TestPromptService_CreateOrUpdatePrompt(t *testing.T) {
	tests := []struct {
		name             string
		promptRepository PromptRepository
		prompt           *Prompt
		want             *Prompt
		wantErr          bool
	}{
		{
			name:             "test case 1",
			promptRepository: &MockPromptRepository{},
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
			service := NewPromptService(tt.promptRepository)
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

			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("CreateOrUpdatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
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
