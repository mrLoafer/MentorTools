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
	Transcription string     `json:"Transcription"`
	Translation   string     `json:"Translation"`
	Description   string     `json:"Description"`
	Synonyms      []string   `json:"Synonyms"`
	Examples      []Examples `json:"Examples"`
}

// Examples представляет пример использования слова
type Examples struct {
	Text        string   `json:"sentence"`
	Keywords    []string `json:"keywords"`
	Translation string   `json:"translation"`
	Context     string   `json:"area"`
}

// Единственное сообщение-инструкция, которое будет храниться в контексте
var systemMessage = Message{
	Role:    "system",
	Content: "Я отправлю тебе слово на английском языке (изучаемое слово) и слова-маркеры для генерации примеров.Ты должен сгенерировать следующую информацию строго в формате JSON без форматирования или дополнительных символов, таких как тройные кавычки:1. Transcription — string: транскрипция изучаемого слова на British English;2. Translation — string: перевод изучаемого слова на русский язык;3. Description — string: объяснение значения изучаемого слова на русском языке;4. Synonyms — []string: список синонимов для изучаемого слова на английском языке;5. Examples — []Examples: список из 5 или более примеров. Структура Examples:1. sentence — string: пример, должен быть максимально простым и коротким, но обязательно содержать изучаемое слово и слова-маркеры;2. keywords — []string: перечень неизвестных, новых, других слов, которые используются в примере (больше-лучше);3. translation — string: перевод примера на русский;4. area — string: область или контекст примера одним словом или выражением."}

// GetWordDetailsFromGPT - Функция для получения данных о слове от ChatGPT
func GetWordDetailsFromGPT(word string, wordsForExamples []string) (string, string, string, []string, []Examples, error) {
	// Формируем новое сообщение с текущим словом
	userMessage := Message{
		Role: "user",
		Content: fmt.Sprintf(`{
			"word": "%s",
			"wordsForExamples": %v,
		}`, word, wordsForExamples),
	}

	// Формируем запрос, включая только инструкцию и текущее сообщение
	requestBody := OpenAIRequest{
		Model:    "gpt-4-turbo",
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
		return WordDetails{}, fmt.Errorf("failed to parse GPT response1: %v", err)
	}
	return wordDetails, nil
}
