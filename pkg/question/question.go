package question

import "time"

// Question is our structure, each question has a body and options
type Question struct {
	// ID will identify our Question
	ID int32 `json:"id"`
	// Body defines an string which will be use to store questions
	Body string `json:"body"`

	// Options each option contains a body which is an string and a correct which is a boolean
	Options []struct {
		Body    string `json:"body,omitempty"`
		Correct bool   `json:"correct,omitempty"`
	} `json:"options,omitempty"`

	// CreatedAt will contains the creation timestamp
	CreatedAt time.Time `json:"created_at,omitempty"`

	// UpdatedAt will contains the last update timestamp
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Questions is an array of Question structs
type Questions []Question
