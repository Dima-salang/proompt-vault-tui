package vault

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreateOrUpdatePrompt_Unit(t *testing.T) {
	repo := NewFakePromptRepository()
	service := NewPromptService(repo)

	tests := []struct {
		name               string // description of this test case
		prompt             *Prompt
		want               *Prompt
		wantErr            bool
		failCreateOrUpdate bool
		failDelete         bool
		failGetByID        bool
		failGetAll         bool
	}{
		{
			name: "Create Prompt Success",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			want: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr:            false,
			failCreateOrUpdate: false,
		},
		{
			name: "Create Prompt Failure",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr:            true,
			failCreateOrUpdate: true,
		},

		{
			name: "Create Prompt Failure - Title is empty",
			prompt: &Prompt{
				Title:         "",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr:            true,
			failCreateOrUpdate: false,
		},
		{
			name: "Create Prompt Failure - Prompt Content is empty",
			prompt: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "",
			},
			wantErr:            true,
			failCreateOrUpdate: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.failCreateOrUpdate = tt.failCreateOrUpdate
			got, gotErr := service.CreateOrUpdatePrompt(tt.prompt)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateOrUpdatePrompt() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateOrUpdatePrompt() succeeded unexpectedly")
			}

			// compare the prompts
			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("CreateOrUpdatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}





func Test_promptService_DeletePrompt(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      int
		wantErr bool
		failDelete bool
	}{
		{
			name: "Delete Prompt Success",
			id: 1,
			wantErr: false,
			failDelete: false,
		},
		{
			name: "Delete Prompt Failure",
			id: 1,
			wantErr: true,
			failDelete: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakePromptRepository()
			service := NewPromptService(repo)

			repo.prompts[tt.id] = &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			}


			repo.failDelete = tt.failDelete
			gotErr := service.DeletePrompt(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DeletePrompt() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DeletePrompt() succeeded unexpectedly")
			}
		})
	}
}

func Test_promptService_GetPromptByID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      int
		want    *Prompt
		wantErr bool
		failGetByID bool
	}{
		{
			name: "Get Prompt by ID Success",
			id: 1,
			want: &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			},
			wantErr: false,
			failGetByID: false,
		},
		{
			name: "Get Prompt by ID Failure",
			id: 1,
			wantErr: true,
			failGetByID: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakePromptRepository()
			service := NewPromptService(repo)

			repo.prompts[tt.id] = &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			}
			repo.failGetByID = tt.failGetByID
			got, gotErr := service.GetPromptByID(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetPromptByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetPromptByID() succeeded unexpectedly")
			}
			if got.Title != tt.want.Title || got.Description != tt.want.Description || got.PromptContent != tt.want.PromptContent {
				t.Errorf("GetPromptByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_promptService_GetAllPrompts(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    []Prompt
		wantErr bool
		failGetAll bool
	}{
		{
			name: "Get All Prompts Success",
			want: []Prompt{
				{
					Title:         "test title",
					Description:   "test description",
					PromptContent: "test prompt content",
				},
			},
			wantErr: false,
			failGetAll: false,
		},
		{
			name: "Get All Prompts Failure",
			wantErr: true,
			failGetAll: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakePromptRepository()
			service := NewPromptService(repo)

			repo.prompts[1] = &Prompt{
				Title:         "test title",
				Description:   "test description",
				PromptContent: "test prompt content",
			}
			repo.failGetAll = tt.failGetAll
			got, gotErr := service.GetAllPrompts()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAllPrompts() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAllPrompts() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllPrompts() = %v, want %v", got, tt.want)
			}

		})
	}
}



type fakePromptRepository struct {
	prompts            map[int]*Prompt
	nextID             int
	failCreateOrUpdate bool
	failDelete         bool
	failGetByID        bool
	failGetAll         bool
}

// fake prompt repository
func NewFakePromptRepository() *fakePromptRepository {
	return &fakePromptRepository{
		prompts: make(map[int]*Prompt),
		nextID:  0,
	}
}

func (repo *fakePromptRepository) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	if prompt.Title == "" {
		return nil, errors.New("title is required")
	}

	if prompt.PromptContent == "" {
		return nil, errors.New("prompt content is required")
	}

	if repo.failCreateOrUpdate {
		return nil, errors.New("failed to create or update prompt")
	}

	if (prompt.ID == 0) {
		prompt.ID = repo.nextID
		repo.nextID++
	}
	repo.prompts[prompt.ID] = prompt
	return prompt, nil
}

func (repo *fakePromptRepository) DeletePrompt(id int) error {
	if repo.failDelete {
		return errors.New("failed to delete prompt")
	}
	delete(repo.prompts, id)
	return nil
}

func (repo *fakePromptRepository) GetPromptByID(id int) (*Prompt, error) {
	if repo.failGetByID {
		return nil, errors.New("prompt not found")
	}

	// check if the prompt exists
	prompt, exists := repo.prompts[id]
	if !exists {
		return nil, errors.New("prompt not found")
	}
	return prompt, nil
}

func (repo *fakePromptRepository) GetAllPrompts() ([]Prompt, error) {
	if repo.failGetAll {
		return nil, errors.New("failed to get all prompts")
	}
	var prompts []Prompt
	for _, prompt := range repo.prompts {
		prompts = append(prompts, *prompt)
	}
	return prompts, nil
}