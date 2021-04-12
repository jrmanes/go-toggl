package question

// Question is our structure, each question has a body and options
type Question struct {
	// ID will identify our Question
	ID int32 `json:"id"`
	// Body defines an string which will be use to store questions
	Body string `json:"body"`

	// Options each option contains a body which is an string and a correct which is a boolean
	Options []struct {
		Body    string `json:"body"`
		Correct bool   `json:"correct"`
	} `json:"options,omitempty"`
}
type Options struct {
	Body    string  `json:"body"`
	Correct []uint8 `json:"correct"`
}

// Questions is an array of Question structs
type Questions []Question
