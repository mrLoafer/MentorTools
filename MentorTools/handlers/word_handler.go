package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"MentorTools/models"
	"MentorTools/services"

	"github.com/jackc/pgx/v4/pgxpool"
)

// GetWordsHandler - обработчик для получения всех слов ученика с дополнительной информацией
func GetWordsHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение userID из токена
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Запрос для получения слов ученика с транскрипцией, переводом и описанием
		rows, err := dbpool.Query(context.Background(), `
            SELECT w.id, w.word, w.transcription, w.translation, s.status
            FROM words w
            JOIN student_words s ON w.id = s.word_id
            WHERE s.student_id = $1`, userID)

		if err != nil {
			http.Error(w, "Failed to retrieve words", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Массив для хранения слов
		var words []map[string]interface{}
		for rows.Next() {
			var id int
			var word, transcription, translation, status string

			// Сканирование данных из запроса
			if err := rows.Scan(&id, &word, &transcription, &translation, &status); err != nil {
				http.Error(w, "Error scanning word data", http.StatusInternalServerError)
				return
			}

			// Добавление слова в массив с дополнительной информацией
			words = append(words, map[string]interface{}{
				"id":            id,
				"word":          word,
				"transcription": transcription,
				"translation":   translation,
				"status":        status,
			})
		}

		// Возвращение JSON с массивом слов
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"words": words})
	}
}

// AddWordHandler - обработчик для добавления нового слова с использованием GPT
// Обработчик для добавления нового слова
func AddWordHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var newWord struct {
			Word string `json:"word"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newWord); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Убираем лишние пробелы и приводим к нижнему регистру
		word := strings.TrimSpace(strings.ToLower(newWord.Word))

		// Проверяем, существует ли уже это слово в базе данных
		var wordID int
		var wordStatus string

		err = dbpool.QueryRow(context.Background(), "SELECT id, status FROM words WHERE word = $1", word).Scan(&wordID, &wordStatus)
		if err == nil {
			// Если слово существует и его статус не "pending"
			if wordStatus != "pending" {
				// Добавляем связь между пользователем и уже существующим словом
				_, err := dbpool.Exec(context.Background(), "INSERT INTO student_words (student_id, word_id, status) VALUES ($1, $2, $3)", userID, wordID, "need to learn")
				if err != nil {
					http.Error(w, "Failed to add word to user", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]bool{"success": true})
				return
			}
			// Если слово в статусе "pending", продолжаем и запрашиваем данные у OpenAI
		}

		// Извлекаем дополнительные слова ученика (рандомно)
		var wordsForExamples []string

		// Получаем 5 рандомных слова, которые ученик уже выучил
		rows, err := dbpool.Query(context.Background(), "SELECT w.word FROM words w JOIN student_words sw ON w.id = sw.word_id WHERE sw.student_id = $1 AND sw.status = 'learned' ORDER BY random() LIMIT 5", userID)
		if err != nil {
			http.Error(w, "Failed to fetch learned words", http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var learnedWord string
			if err := rows.Scan(&learnedWord); err != nil {
				http.Error(w, "Failed to read learned word", http.StatusInternalServerError)
				return
			}
			wordsForExamples = append(wordsForExamples, learnedWord)
		}

		// Получаем 2 рандомных слов, которые нужно выучить
		rowss, err := dbpool.Query(context.Background(), "SELECT w.word FROM words w JOIN student_words sw ON w.id = sw.word_id WHERE sw.student_id = $1 AND sw.status = 'need to learn' ORDER BY random() LIMIT 2", userID)
		if err != nil {
			http.Error(w, "Failed to fetch upcoming words!", http.StatusInternalServerError)
			return
		}
		for rowss.Next() {
			var upcomingWord string
			if err := rowss.Scan(&upcomingWord); err != nil {
				http.Error(w, "Failed to read upcoming word", http.StatusInternalServerError)
				return
			}
			wordsForExamples = append(wordsForExamples, upcomingWord)
		}

		// Запрашиваем данные у OpenAI, включая дополнительные слова
		transcription, translation, description, synonyms, examples, err := services.GetWordDetailsFromGPT(word, wordsForExamples)
		if err != nil {
			http.Error(w, "Failed to fetch word details from GPT", http.StatusInternalServerError)
			return
		}

		if wordID == 0 {
			// Если слово не существует в БД, добавляем его
			err = dbpool.QueryRow(context.Background(), "INSERT INTO words (word, transcription, translation, definition, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				word, transcription, translation, description, "completed").Scan(&wordID)
			if err != nil {
				http.Error(w, "Failed to add word", http.StatusInternalServerError)
				return
			}
		} else {
			// Если слово в статусе "pending", обновляем его информацию
			_, err = dbpool.Exec(context.Background(), "UPDATE words SET transcription = $1, translation = $2, definition = $3, status = 'completed' WHERE id = $4",
				transcription, translation, description, wordID)
			if err != nil {
				http.Error(w, "Failed to update word", http.StatusInternalServerError)
				return
			}
		}

		// Добавляем связь с пользователем
		_, err = dbpool.Exec(context.Background(), "INSERT INTO student_words (student_id, word_id, status) VALUES ($1, $2, $3)", userID, wordID, "need to learn")
		if err != nil {
			http.Error(w, "Failed to add word to user", http.StatusInternalServerError)
			return
		}

		// Сохраняем синонимы (если есть)
		for _, synonym := range synonyms {
			var synonymID int
			err := dbpool.QueryRow(context.Background(), "SELECT id FROM words WHERE word = $1", synonym).Scan(&synonymID)
			if err != nil {
				err = dbpool.QueryRow(context.Background(), "INSERT INTO words (word, status) VALUES ($1, $2) RETURNING id", synonym, "pending").Scan(&synonymID)
				if err != nil {
					http.Error(w, "Failed to save synonym", http.StatusInternalServerError)
					return
				}
			}
			_, err = dbpool.Exec(context.Background(), "INSERT INTO word_links (word_id, linked_word_id) VALUES ($1, $2)", wordID, synonymID)
			if err != nil {
				http.Error(w, "Failed to link synonym", http.StatusInternalServerError)
				return
			}
		}

		// Сохраняем примеры и ключевые слова
		for _, example := range examples {
			var exampleID int
			err = dbpool.QueryRow(context.Background(), "INSERT INTO examples (example, context, translation) VALUES ($1, $2, $3) RETURNING id", example.Text, example.Context, example.Translation).Scan(&exampleID)
			if err != nil {
				http.Error(w, "Failed to save example", http.StatusInternalServerError)
				return
			}

			_, err = dbpool.Exec(context.Background(), "INSERT INTO word_example (word_id, example_id) VALUES ($1, $2)", wordID, exampleID)
			if err != nil {
				http.Error(w, "Failed to link word and example", http.StatusInternalServerError)
				return
			}

			for _, keyword := range example.Keywords {
				var keywordID int
				err := dbpool.QueryRow(context.Background(), "SELECT id FROM words WHERE word = $1", keyword).Scan(&keywordID)

				if err != nil {
					// Если слово не найдено, добавляем его с статусом 'pending'
					err = dbpool.QueryRow(context.Background(), "INSERT INTO words (word, status) VALUES ($1, $2) RETURNING id", keyword, "pending").Scan(&keywordID)

					if err != nil {
						http.Error(w, "Failed to save keyword", http.StatusInternalServerError)
						return
					}
				}

				// Проверяем, существует ли уже связь между словом и примером
				var exists bool
				err = dbpool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM word_example WHERE word_id = $1 AND example_id = $2)", keywordID, exampleID).Scan(&exists)
				if err != nil {
					http.Error(w, "Failed to check if word and example are linked", http.StatusInternalServerError)
					return
				}

				// Если связь уже существует, просто продолжаем цикл
				if exists {
					continue // Продолжаем цикл, чтобы обработать остальные ключевые слова
				}

				// Если связи нет, создаем ее
				_, err = dbpool.Exec(context.Background(), "INSERT INTO word_example (word_id, example_id) VALUES ($1, $2)", keywordID, exampleID)
				if err != nil {
					http.Error(w, "Failed to link keyword and example", http.StatusInternalServerError)
					return
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// UpdateWordStatusHandler - обработчик для обновления статуса слова
func UpdateWordStatusHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusUpdate struct {
			WordID int    `json:"wordId"`
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = dbpool.Exec(context.Background(), `
            UPDATE student_words
            SET status = $1
            WHERE word_id = $2 AND student_id = $3`, statusUpdate.Status, statusUpdate.WordID, userID)

		if err != nil {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// GetWordDetailsHandler - обработчик для получения деталей по слову
func GetWordDetailsHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wordID := r.URL.Query().Get("id")

		var word models.WordDetails
		err := dbpool.QueryRow(context.Background(), `
            SELECT w.word, w.transcription, w.translation, w.definition
            FROM words w
            WHERE w.id = $1`, wordID).Scan(&word.Word, &word.Transcription, &word.Translation, &word.Description)

		if err != nil {
			http.Error(w, "Word not found", http.StatusNotFound)
			return
		}

		// Получаем синонимы
		rows, err := dbpool.Query(context.Background(), "SELECT linked_word_id FROM word_links WHERE word_id = $1", wordID)
		if err != nil {
			http.Error(w, "Failed to fetch synonyms", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var synonyms []string
		for rows.Next() {
			var linkedWordID int
			if err := rows.Scan(&linkedWordID); err != nil {
				http.Error(w, "Error reading synonym ID", http.StatusInternalServerError)
				return
			}

			// Получаем отдельное соединение из пула
			var synonym string
			err = dbpool.QueryRow(context.Background(), "SELECT word FROM words WHERE id = $1", linkedWordID).Scan(&synonym)

			// Если возникает ошибка, логируем её и продолжаем с другими синонимами
			if err != nil {
				continue // Пропускаем синоним в случае ошибки
			}

			synonyms = append(synonyms, synonym)
		}

		// Устанавливаем синонимы для слова
		word.Synonyms = synonyms

		// Получаем примеры
		exampleRows, err := dbpool.Query(context.Background(), `
    SELECT e.example, e.translation
    FROM examples e
    JOIN word_example we ON e.id = we.example_id
    WHERE we.word_id = $1`, wordID)

		if err != nil {
			http.Error(w, "Failed to fetch examples", http.StatusInternalServerError)
			return
		}
		defer exampleRows.Close()

		var examples []models.Example
		for exampleRows.Next() {
			var sentence sql.NullString
			var translation sql.NullString

			if err := exampleRows.Scan(&sentence, &translation); err != nil {
				http.Error(w, "Error reading example", http.StatusInternalServerError)
				return
			}

			var example models.Example

			// Проверяем и присваиваем значение, если оно не NULL
			if sentence.Valid {
				example.Sentence = sentence.String
			} else {
				example.Sentence = "" // Или установите какое-то значение по умолчанию
			}

			if translation.Valid {
				example.Translation = translation.String
			} else {
				example.Translation = "" // Или установите какое-то значение по умолчанию
			}

			examples = append(examples, example)
		}
		word.Examples = examples

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(word)
	}
}
