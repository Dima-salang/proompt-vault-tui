package vault

import (
	"github.com/atotto/clipboard"
	"github.com/sahilm/fuzzy"
)

func SearchPrompts(prompts Prompts, query string) fuzzy.Matches {
	// search for the query in the prompts array
	results := fuzzy.FindFrom(query, prompts)

	// return the results
	return results
}

func CopyToClipboard(prompt *Prompt) error {
	// copy the text to the clipboard
	err := clipboard.WriteAll(prompt.PromptContent)
	if err != nil {
		return err
	}
	return nil
}
