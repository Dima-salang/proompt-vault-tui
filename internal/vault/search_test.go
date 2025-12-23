package vault_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Dima-salang/proompt-vault-tui/internal/vault"
	"github.com/sahilm/fuzzy"
)

func TestSearchPrompts(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prompts vault.Prompts
		query   string
		want    fuzzy.Matches
	}{
		{
			name: "Fuzzy Search for 'test' in prompts array",
			prompts: vault.Prompts{
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
			got := vault.SearchPrompts(tt.prompts, tt.query)

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
		prompt  *vault.Prompt
		wantErr bool
		want    string
	}{
		{
			name: "Copy to clipboard success",
			prompt: &vault.Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: false,
			want:    "test prompt content",
		},
		{
			name: "Copy to clipboard failure",
			prompt: &vault.Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: true,
			want:    "test prompt content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipboard := &FakeClipBoard{
				clipboard: tt.want,
				wantErr:   tt.wantErr,
			}
			gotErr := clipboard.CopyToClipboard(tt.prompt.PromptContent)

			// writing to the clipboard
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CopyToClipboard() failed: %v", gotErr)
				}
				return
			}

			if clipboard.clipboard != tt.want {
				t.Errorf("CopyToClipboard() = %v, want %v", clipboard.clipboard, tt.want)
			}

			if tt.wantErr {
				t.Fatal("CopyToClipboard() succeeded unexpectedly")
			}
		})
	}
}

// fake clipboard implementation for unit tests
// since github actions does not have any clipboard functionality as it is headless
type FakeClipBoard struct {
	clipboard string
	wantErr   bool
}

func (f *FakeClipBoard) CopyToClipboard(textToCopy string) error {
	if f.wantErr {
		return errors.New("failed to write to clipboard")
	}
	f.clipboard = textToCopy
	return nil
}

func (f *FakeClipBoard) ReadAll() (string, error) {
	if f.wantErr {
		return "", errors.New("failed to read from clipboard")
	}
	return f.clipboard, nil
}
