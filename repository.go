package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"log/slog"
	"github.com/boltdb/bolt"
)

var logger *slog.Logger

// get bolt db connection
func getDB() (*bolt.DB, error) {
	// open the database
	db, err := bolt.Open("prompt.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}



type PromptRepository interface {
	CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error)
	DeletePrompt(id int) error
	GetPromptByID(id int) (*Prompt, error)
	GetAllPrompts() ([]Prompt, error)
}


// creates or updates an individual prompt
func CreateOrUpdatePrompt(prompt *Prompt) (*Prompt, error) {
	// get the prompt bucket
	db, err := getDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	// write the prompt to the bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			logger.Error("failed to create bucket", "error", err)
			return err
		}

		// create a unique key for the prompt
		// replaces the id if it already exists
		id, _ := bucket.NextSequence()
		prompt.ID = int(id)

		// encode the prompt
		encodedPrompt, err := json.Marshal(prompt)
		if err != nil {
			logger.Error("failed to encode prompt", "error", err)
			return err
		}

		// write the prompt to the bucket
		key := itob(uint64(id))
		err = bucket.Put(key, encodedPrompt)
		if err != nil {
			logger.Error("failed to write prompt to bucket", "error", err)
			return err
		}

		return nil
	})

	return prompt, err
}


// delete the prompt
func DeletePrompt(id int) error {
	// get the prompt bucket
	db, err := getDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		log.Fatal(err)
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			logger.Error("failed to create bucket", "error", err)
			return err
		}

		// delete the prompt
		key := itob(uint64(id))
		err = bucket.Delete(key)
		if err != nil {
			logger.Error("failed to delete prompt", "error", err)
			return err
		}

		return nil
	})

	return err
}


// get specific prompt details by id
func GetPromptByID(id int) (*Prompt, error) {
	// get the prompt bucket
	db, err := getDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	// create empty prompt struct
	prompt := &Prompt{}

	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			logger.Error("failed to create bucket", "error", err)
			return err
		}

		// get the prompt
		key := itob(uint64(id))
		value := bucket.Get(key)
		if value == nil {
			logger.Error("prompt not found", "id", id)
			return nil
		}

		// decode the prompt
		err = json.Unmarshal(value, prompt)
		if err != nil {
			logger.Error("failed to decode prompt", "error", err)
			return err
		}

		return nil
	})
	



	return prompt, err
}



// get all prompts
func GetAllPrompts() ([]Prompt, error) {
	// get the prompt bucket
	db, err := getDB()
	if err != nil {
		logger.Error("failed to open database", "error", err)
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	// create empty prompt struct
	prompts := []Prompt{}

	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			return err
		}

		// get all prompts
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			prompt := &Prompt{}
			err = json.Unmarshal(v, prompt)
			if err != nil {
				logger.Error("failed to decode prompt", "error", err)
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


	

type ProjectRepository interface {
	CreateProject(project *Project) error
	UpdateProject(project *Project) error
	DeleteProject(id int) error
	GetProjectByID(id int) (*Project, error)
	GetAllProjects() ([]Project, error)
}




	
