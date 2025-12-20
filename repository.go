package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"github.com/boltdb/bolt"
)

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
	CreateOrUpdatePrompt(prompt *Prompt) (Prompt, error)
	DeletePrompt(id int) error
	GetPromptByID(id int) (*Prompt, error)
	GetAllPrompts() ([]Prompt, error)
}


// creates or updates an individual prompt
func CreateOrUpdatePrompt(prompt *Prompt) (Prompt, error) {
	// get the prompt bucket
	db, err := getDB()
	if err != nil {
		log.Println(err)
		return Prompt{}, err
	}
	defer db.Close()

	// write the prompt to the bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("prompts"))
		if err != nil {
			return err
		}

		// create a unique key for the prompt
		// replaces the id if it already exists
		id, _ := bucket.NextSequence()
		prompt.ID = int(id)

		// encode the prompt
		encodedPrompt, err := json.Marshal(prompt)
		if err != nil {
			return err
		}

		// write the prompt to the bucket
		key := itob(uint64(id))
		err = bucket.Put(key, encodedPrompt)
		if err != nil {
			return err
		}

		return nil
	})

	return *prompt, err
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




	
