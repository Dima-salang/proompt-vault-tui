package main

import "errors"

type PromptService struct {
	promptRepository PromptRepository
}

// creates a new prompt service
func NewPromptService(promptRepository PromptRepository) *PromptService {
	return &PromptService{
		promptRepository: promptRepository,
	}
}

// creates or updates an individual prompt
func (service *PromptService) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	// validate the prompt
	if prompt.Title == "" {
		return nil, errors.New("title is required")
	}
	if prompt.PromptContent == "" {
		return nil, errors.New("prompt content is required")
	}
	return service.promptRepository.CreateOrUpdatePrompt(prompt)
}


func (service *PromptService) DeletePrompt(id int) error {
	return service.promptRepository.DeletePrompt(id)
}


func (service *PromptService) GetPromptByID(id int) (*Prompt, error) {
	return service.promptRepository.GetPromptByID(id)
}


func (service *PromptService) GetAllPrompts() ([]Prompt, error) {
	return service.promptRepository.GetAllPrompts()
}
