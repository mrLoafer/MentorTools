package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Структуры для работы с OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Функция для получения данных о слове от ChatGPT
func GetWordDetailsFromGPT(word string) (string, string, string, error) {
	requestBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "Ты — помощник по изучению английского языка. Предоставь транскрипцию, описание и примеры для данного слова.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Слово: %s", word),
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", "", fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var openAIResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return "", "", "", fmt.Errorf("error decoding response: %v", err)
	}

	// Проверка на наличие ответа
	if len(openAIResponse.Choices) == 0 {
		return "", "", "", fmt.Errorf("no choices returned from GPT")
	}

	// Получаем содержимое ответа
	content := openAIResponse.Choices[0].Message.Content

	// Парсим содержимое ответа
	transcription, description, examples := parseGPTResponse(content)

	return transcription, description, examples, nil
}

// Функция для парсинга ответа от ChatGPT
func parseGPTResponse(response string) (string, string, string) {
	// Ищем блоки "Transcription:", "Description:" и "Examples:"
	transcription := extractBetween(response, "Transcription:", "Description:")
	description := extractBetween(response, "Description:", "Examples:")
	examples := extractExamples(response)

	return strings.TrimSpace(transcription), strings.TrimSpace(description), strings.TrimSpace(examples)
}

// Вспомогательная функция для извлечения текста между двумя подстроками
func extractBetween(value string, start string, end string) string {
	s := strings.Index(value, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(value[s:], end)
	if e == -1 {
		return value[s:]
	}
	return value[s : s+e]
}

// Вспомогательная функция для извлечения примеров
func extractExamples(value string) string {
	examplesStart := strings.Index(value, "Examples:")
	if examplesStart == -1 {
		return ""
	}
	examples := value[examplesStart+len("Examples:"):]

	// Преобразуем примеры в читаемый формат
	examples = strings.ReplaceAll(examples, "\n", " ")
	examples = strings.TrimSpace(examples)

	return examples
}
