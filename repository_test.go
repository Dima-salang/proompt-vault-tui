package main_test

import (
	"testing"

	main "github.com/Dima-salang/proompt-vault-tui"
)

func TestCreateOrUpdatePrompt(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prompt  *main.Prompt
		want    *main.Prompt
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case 1",
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
			db, err := main.getDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			repo := main.NewPromptRepository(db)
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
			// TODO: update the condition below to compare got with tt.want.

			// compare the prompts
			// does not compare the id as it is auto generated from the db
			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("CreateOrUpdatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
