package vault

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/boltdb/bolt"
)

type PromptRepository interface {
	CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error)
	DeletePrompt(id int) error
	GetPromptByID(id int) (*Prompt, error)
	GetAllPrompts() ([]Prompt, error)
}

type promptRepository struct {
	db     *bolt.DB
	logger *slog.Logger
}

// creates a new prompt repository
func NewPromptRepository(db *bolt.DB, logger *slog.Logger) PromptRepository {
	return &promptRepository{db: db, logger: logger}
}

// creates or updates an individual prompt
func (repo *promptRepository) CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	// get the prompt bucket
	db := repo.db

	// write the prompt to the bucket
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			repo.logger.Error("failed to create bucket", "error", err)
			return err
		}

		// create a unique key for the prompt
		// replaces the id if it already exists
		id, _ := bucket.NextSequence()
		prompt.ID = int(id)

		// encode the prompt
		encodedPrompt, err := json.Marshal(prompt)
		if err != nil {
			repo.logger.Error("failed to encode prompt", "error", err)
			return err
		}

		// write the prompt to the bucket
		key := itob(uint64(id))
		err = bucket.Put(key, encodedPrompt)
		if err != nil {
			repo.logger.Error("failed to write prompt to bucket", "error", err)
			return err
		}

		return nil
	})

	return prompt, err
}

// delete the prompt
func (repo *promptRepository) DeletePrompt(id int) error {
	// get the prompt bucket
	db := repo.db

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			repo.logger.Error("failed to create bucket", "error", err)
			return err
		}

		// delete the prompt
		key := itob(uint64(id))
		err = bucket.Delete(key)
		if err != nil {
			repo.logger.Error("failed to delete prompt", "error", err)
			return err
		}

		return nil
	})

	return err
}

// get specific prompt details by id
func (repo *promptRepository) GetPromptByID(id int) (*Prompt, error) {
	// get the prompt bucket
	db := repo.db

	// create empty prompt struct
	prompt := &Prompt{}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("prompts"))
		if bucket == nil {
			repo.logger.Error("bucket not found")
			return errors.New("prompt not found")
		}

		// get the prompt
		key := itob(uint64(id))
		value := bucket.Get(key)
		if value == nil {
			repo.logger.Error("prompt not found", "id", id)
			return errors.New("prompt not found")
		}

		// decode the prompt
		err := json.Unmarshal(value, prompt)
		if err != nil {
			repo.logger.Error("failed to decode prompt", "error", err)
			return err
		}

		return nil
	})

	return prompt, err
}

// get all prompts
func (repo *promptRepository) GetAllPrompts() ([]Prompt, error) {
	// get the prompt bucket
	db := repo.db

	// create empty prompt struct
	prompts := []Prompt{}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("prompts"))
		if bucket == nil {
			return nil // No prompts yet
		}

		// get all prompts
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			prompt := &Prompt{}
			err := json.Unmarshal(v, prompt)
			if err != nil {
				repo.logger.Error("failed to decode prompt", "error", err)
				return err
			}
			prompts = append(prompts, *prompt)
		}

		return nil
	})

	return prompts, err
}

// helper function to convert uint64 to []byte
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
