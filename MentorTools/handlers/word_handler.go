package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"MentorTools/services"

	"github.com/jackc/pgx/v4"
)

// GetWordsHandler - обработчик для получения всех слов ученика
func GetWordsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение userID из токена
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(context.Background(), `
            SELECT w.id, w.word, s.status
            FROM words w
            JOIN student_words s ON w.id = s.word_id
            WHERE s.student_id = $1`, userID)

		if err != nil {
			http.Error(w, "Failed to retrieve words", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var words []map[string]interface{}
		for rows.Next() {
			var word string
			var id int
			var status string
			if err := rows.Scan(&id, &word, &status); err != nil {
				http.Error(w, "Error scanning word", http.StatusInternalServerError)
				return
			}

			words = append(words, map[string]interface{}{
				"id":     id,
				"word":   word,
				"status": status,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"words": words})
	}
}

// AddWordHandler - обработчик для добавления нового слова
func AddWordHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newWord struct {
			Word          string `json:"word"`
			Transcription string `json:"transcription"`
			Definition    string `json:"definition"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newWord); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Очистка лишних пробелов и перевод слова в нижний регистр для проверки уникальности
		trimmedWord := strings.ToLower(strings.TrimSpace(newWord.Word))

		// Проверяем наличие слова в базе данных
		var wordID int
		err := db.QueryRow(context.Background(), `
			SELECT id FROM words WHERE LOWER(TRIM(word)) = $1`, trimmedWord).Scan(&wordID)

		if err == pgx.ErrNoRows {
			// Если слово не найдено, создаем новое слово
			err = db.QueryRow(context.Background(), `
				INSERT INTO words (word, transcription, definition)
				VALUES ($1, $2, $3) RETURNING id`, trimmedWord, newWord.Transcription, newWord.Definition).Scan(&wordID)

			if err != nil {
				http.Error(w, "Failed to add word", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, "Failed to check word", http.StatusInternalServerError)
			return
		}

		// Получение ID ученика из токена (должен быть реализован метод для извлечения ID)
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверяем, есть ли уже связь между учеником и словом
		var existingLink int
		err = db.QueryRow(context.Background(), `
			SELECT 1 FROM student_words WHERE student_id = $1 AND word_id = $2`, userID, wordID).Scan(&existingLink)

		if err == pgx.ErrNoRows {
			// Если связи нет, создаем новую связь между учеником и словом
			_, err = db.Exec(context.Background(), `
				INSERT INTO student_words (student_id, word_id, status)
				VALUES ($1, $2, 'need to learn')`, userID, wordID)

			if err != nil {
				http.Error(w, "Failed to link word to student", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, "Failed to check student-word link", http.StatusInternalServerError)
			return
		}

		// Возвращаем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// UpdateWordStatusHandler - обработчик для обновления статуса слова
func UpdateWordStatusHandler(db *pgx.Conn) http.HandlerFunc {
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

		_, err = db.Exec(context.Background(), `
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
func GetWordDetailsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wordID := r.URL.Query().Get("id")
		if wordID == "" {
			http.Error(w, "Missing word ID", http.StatusBadRequest)
			return
		}

		var word struct {
			Word          string `json:"word"`
			Transcription string `json:"transcription"`
			Definition    string `json:"definition"`
		}

		err := db.QueryRow(context.Background(), `
            SELECT word, transcription, definition 
            FROM words 
            WHERE id = $1`, wordID).Scan(&word.Word, &word.Transcription, &word.Definition)

		if err != nil {
			http.Error(w, "Failed to retrieve word details", http.StatusInternalServerError)
			return
		}

		// Получение примеров для слова
		rows, err := db.Query(context.Background(), `
            SELECT e.example
            FROM examples e
            JOIN word_example we ON e.id = we.example_id
            WHERE we.word_id = $1`, wordID)

		if err != nil {
			http.Error(w, "Failed to retrieve examples", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var examples []string
		for rows.Next() {
			var example string
			if err := rows.Scan(&example); err != nil {
				http.Error(w, "Error reading examples", http.StatusInternalServerError)
				return
			}
			examples = append(examples, example)
		}

		response := map[string]interface{}{
			"word":          word.Word,
			"transcription": word.Transcription,
			"definition":    word.Definition,
			"examples":      examples,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
