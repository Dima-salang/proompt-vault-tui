package vault

import "errors"

type PromptService interface {
	CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error)
	DeletePrompt(id int) error
	GetPromptByID(id int) (*Prompt, error)
	GetAllPrompts() ([]Prompt, error)
}

type promptService struct {
	promptRepository PromptRepository
}

// creates a new prompt service
func NewPromptService(promptRepository PromptRepository) PromptService {
	return &promptService{
		promptRepository: promptRepository,
	}
}

/* 
	Creates or updates an individual prompt.
	It returns the created or updated prompt and an error if any.
*/
func (service *promptService) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	// validate the prompt
	if prompt.Title == "" {
		return nil, errors.New("title is required")
	}
	if prompt.PromptContent == "" {
		return nil, errors.New("prompt content is required")
	}
	return service.promptRepository.CreateOrUpdatePrompt(prompt)
}


// Deletes a specific prompt by ID.
func (service *promptService) DeletePrompt(id int) error {
	return service.promptRepository.DeletePrompt(id)
}


// Gets a specific prompt by ID.
func (service *promptService) GetPromptByID(id int) (*Prompt, error) {
	return service.promptRepository.GetPromptByID(id)
}


// Gets all prompts.
// Primarily used for the list view and initial loading of the app
func (service *promptService) GetAllPrompts() ([]Prompt, error) {
	return service.promptRepository.GetAllPrompts()
}
