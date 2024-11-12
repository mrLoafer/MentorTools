CREATE TABLE IF NOT EXISTS user_links (
                                          id SERIAL PRIMARY KEY,
                                          teacher_id INT NOT NULL,
                                          student_id INT NOT NULL,
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                          FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
                                          FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);