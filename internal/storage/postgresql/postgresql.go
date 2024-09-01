package postgresql

import (
	"database/sql"
	"fmt"
	task "url-shorter/models/task"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB // поле конекта к бд
}

func New(storagePath string) (*Storage, error) {
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	const op = "storage.postgresql.New"

	// connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	// db, err := sql.Open("postgres", connStr)

	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		due_date TIMESTAMP WITH TIME ZONE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) SaveTask(task *task.Task) error {
	const op = "storage.postgresql.SaveTask"

	query := "INSERT INTO tasks (title, description, due_date) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRow(query, task.Title, task.Description, task.DueDate).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

//{  Post должен вернуть ответ
// "id": "int",
// "title": "string",
// "description": "string",
// "due_date": "string (RFC3339 format)",
// "created_at": "string (RFC3339 format)", "updated_at": "string (RFC3339 format)" }

func (s *Storage) GetTask(id int) (*task.Task, error) {
	const op = "storage.postgresql.GetTask"

	getTask := &task.Task{}
	query := `SELECT * FROM tasks WHERE id=$1`
	err := s.db.QueryRow(query, id).Scan(&getTask.ID, &getTask.Title, &getTask.Description, &getTask.DueDate, &getTask.CreatedAt, &getTask.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: task with ID %d not found", op, id)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return getTask, nil

}

func (s *Storage) Close() error {
	return s.db.Close()
}
