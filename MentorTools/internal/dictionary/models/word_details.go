package models

// WordDetails содержит полную информацию о слове, включая синонимы и примеры
type WordDetails struct {
	Word          string    `json:"word"`
	Transcription string    `json:"transcription"`
	Translation   string    `json:"translation"`
	Description   string    `json:"description"`
	Synonyms      []string  `json:"synonyms"`
	Examples      []Example `json:"examples"`
}

// Example содержит пример с предложением и переводом
type Example struct {
	Sentence    string `json:"sentence"`
	Translation string `json:"translation"`
}
