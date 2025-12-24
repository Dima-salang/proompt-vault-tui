package vault

import (
	"github.com/atotto/clipboard"
	"github.com/sahilm/fuzzy"
)

// Searches for prompts that match the query returns the fuzzy matches.
// It uses the source from prompt.go since FindFrom expects a Source type.
// Uses the fuzzy lib from sahilm.
func SearchPrompts(prompts Prompts, query string) fuzzy.Matches {
	// search for the query in the prompts array
	results := fuzzy.FindFrom(query, prompts)

	// return the results
	return results
}

// Copies the prompt content to the clipboard.
// Requires xsel or xclip to be installed for Linux according to the atotto/clipboard docs.
func CopyToClipboard(prompt *Prompt) error {
	// copy the text to the clipboard
	err := clipboard.WriteAll(prompt.PromptContent)
	if err != nil {
		return err
	}
	return nil
}
