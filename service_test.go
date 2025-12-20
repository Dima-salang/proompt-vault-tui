package main_test

import (
	"testing"

	main "github.com/Dima-salang/proompt-vault-tui"
)

func TestPromptService_CreateOrUpdatePrompt(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		promptRepository main.PromptRepository
		// Named input parameters for target function.
		prompt  *main.Prompt
		want    *main.Prompt

		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:             "test case 1",
			promptRepository: main.NewPromptRepository(),
			prompt: &main.Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			want: &main.Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := main.NewPromptService(tt.promptRepository)
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
			// TODO: update the condition below to compare got with tt.want.
			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("CreateOrUpdatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
