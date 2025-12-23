package main

import "time"

// Project struct
type Project struct {
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Prompt struct
type Prompt struct {
	ID            int
	ProjectID     int
	Title         string
	Description   string
	PromptContent string
	Tags          []string
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
