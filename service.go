package main

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
	return service.promptRepository.CreateOrUpdatePrompt(prompt)
}