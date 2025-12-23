package main_test

import (
	"reflect"
	"testing"

	main "github.com/Dima-salang/proompt-vault-tui"
	"github.com/atotto/clipboard"
	"github.com/sahilm/fuzzy"
)

func TestSearchPrompts(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prompts main.Prompts
		query   string
		want    fuzzy.Matches
	}{
		{
			name: "Fuzzy Search for 'test' in prompts array",
			prompts: main.Prompts{
				{
					Title:         "test title",
					Description:   "test description",
					PromptContent: "test prompt content",
				},
			},
			query: "test",
			want: fuzzy.Matches{
				// since the score is a rank from the algorithm, we need to run it first to know the exact score
				// funnily it's 69...
				{
					Str:            "test title",
					Index:          0,
					MatchedIndexes: []int{0, 1, 2, 3},
					Score:          69,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := main.SearchPrompts(tt.prompts, tt.query)

			// compare the results 
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchPrompts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyToClipboard(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prompt  *main.Prompt
		wantErr bool
		want    string
	}{
		{
			name: "Copy to clipboard success",
			prompt: &main.Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: false,
			want:    "test prompt content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := main.CopyToClipboard(tt.prompt)

			// we also read from the clipboard since we want to verify
			// whether it truly wrote to the clipboard
			got, err := clipboard.ReadAll()


			// writing to the clipboard
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CopyToClipboard() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CopyToClipboard() succeeded unexpectedly")
			}

			// reading from the clipboard
			if err != nil {
				t.Errorf("CopyToClipboard() failed: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyToClipboard() = %v, want %v", got, tt.want)
			}
		})
	}
}
