package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Определение структуры запроса к OpenAI
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message представляет отдельное сообщение в чате с OpenAI
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse представляет ответ от OpenAI
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice представляет выбор в ответе от OpenAI
type Choice struct {
	Message Message `json:"message"`
}

// WordDetails структура для транскрипции, перевода, примеров и описания
type WordDetails struct {
	Transcription string    `json:"Transcription"`
	Translation   string    `json:"Translation"`
	Description   string    `json:"Description"`
	Synonyms      []string  `json:"Synonyms"`
	Examples      []Example `json:"Examples"`
}

// Example представляет пример использования слова
type Example struct {
	Text        string   `json:"sentence"`
	Keywords    []string `json:"keywords"`
	Translation string   `json:"translation"`
	Context     string   `json:"area"`
}

// Единственное сообщение-инструкция, которое будет храниться в контексте
var systemMessage = Message{
	Role:    "system",
	Content: "Ты — помощник по изучению английского языка. Я буду отправлять тебе слово на английском языке и несколько дополнительных слов. Ты должен отвечать в формате JSON. Ответ должен быть в формате JSON с полями: 1. Transcription (транскрипция) — string: транскрипция слова на British English; 2. Translation (перевод) — string: перевод слова на русский язык; 3. Description (описание) — string: объяснение значения слова на русском языке; 4. Synonyms (синонимы) — []string: список синонимов; 5. Examples (примеры) — []Example: список из 5 примеров. В каждом примере обязательно должны использоваться все или некоторые из дополнительных слов (learnedWords и upcomingWords) из запроса. Для каждого примера нужно указать: sentence — string: пример использования слова с дополнительными словами; keywords — []string: новые слова в примере; translation — string: перевод примера на русский; area — string: область или контекст примера."}

// Функция для получения данных о слове от ChatGPT
func GetWordDetailsFromGPT(word string, learnedWords []string, upcomingWords string) (string, string, string, []string, []Example, error) {
	// Формируем новое сообщение с текущим словом
	userMessage := Message{
		Role: "user",
		Content: fmt.Sprintf(`{
			"word": "%s",
			"learnedWords": %v,
			"upcomingWords": %v
		}`, word, learnedWords, upcomingWords),
	}

	// Формируем запрос, включая только инструкцию и текущее сообщение
	requestBody := OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []Message{systemMessage, userMessage}, // Инструкция + текущее сообщение
	}

	// Конвертируем запрос в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("error marshalling request: %v", err)
	}

	fmt.Printf("Request body to GPT: %s\n", jsonData) // Логируем тело запроса

	// Создаем HTTP-запрос к OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("error creating request: %v", err)
	}

	// Логируем URL и заголовки запроса
	fmt.Printf("Request URL: %s\n", req.URL)
	fmt.Printf("Request headers: %v\n", req.Header)

	// Устанавливаем необходимые заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Логируем статус ответа
	fmt.Println("GPT Response Status:", resp.StatusCode)

	// Обрабатываем ответ
	var openAIResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return "", "", "", nil, nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Проверяем, что ответ содержит хотя бы один выбор
	if len(openAIResponse.Choices) == 0 {
		return "", "", "", nil, nil, fmt.Errorf("no choices returned from GPT")
	}

	// Логируем полученный ответ от GPT
	content := openAIResponse.Choices[0].Message.Content
	fmt.Println("GPT Response:", content)

	// Парсим содержимое ответа
	wordDetails, err := parseGPTResponse(content)
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("error parsing GPT response: %v", err)
	}

	return wordDetails.Transcription, wordDetails.Translation, wordDetails.Description, wordDetails.Synonyms, wordDetails.Examples, nil
}

// Функция для парсинга ответа в нужную структуру
func parseGPTResponse(content string) (WordDetails, error) {
	var wordDetails WordDetails

	// Пытаемся распарсить JSON-ответ
	err := json.Unmarshal([]byte(content), &wordDetails)
	if err != nil {
		return WordDetails{}, fmt.Errorf("failed to parse GPT response: %v", err)
	}
	return wordDetails, nil
}
