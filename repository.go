package main

type PromptRepository interface {
	CreatePrompt(prompt *Prompt) error
	UpdatePrompt(prompt *Prompt) error
	DeletePrompt(id int) error
	GetPromptByID(id int) (*Prompt, error)
	GetAllPrompts() ([]Prompt, error)
}

type ProjectRepository interface {
	CreateProject(project *Project) error
	UpdateProject(project *Project) error
	DeleteProject(id int) error
	GetProjectByID(id int) (*Project, error)
	GetAllProjects() ([]Project, error)
}




	
