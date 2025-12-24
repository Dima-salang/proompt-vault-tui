package vault

import "time"

// Prompt struct
type Prompt struct {
	ID            int
	Title         string
	Description   string
	PromptContent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// prompts array for fuzzy search
type Prompts []Prompt

func (p Prompts) Len() int {
	return len(p)
}

func (p Prompts) String(i int) string {
	return p[i].Title
}
