package main

import (
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
