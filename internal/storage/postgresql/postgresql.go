package postgresql

import (
	"database/sql"
	"fmt"
	task "to-do-list/models/task"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB // поле конекта к бд
}

func New(storagePath string) (*Database, error) {
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	const op = "storage.postgresql.New"

	// connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	// db, err := sql.Open("postgres", connStr)

	db, err := sql.Open("postgres", "user=postgrest dbname=db password=postgrest sslmode=disable")

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

	return &Database{
		db: db,
	}, nil

}

func (s *Database) InsertTaskRecord(task *task.Task) (*task.Task, error) {
	const op = "storage.postgresql.InsertTaskRecord"

	query := "INSERT INTO tasks (title, description, due_date) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRow(query, task.Title, task.Description, task.DueDate).Scan(&task.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return s.QueryTaskByID(task.ID)
}

func (s *Database) QueryTaskByID(id int) (*task.Task, error) {
	const op = "storage.postgresql.QueryTaskByID"

	task := &task.Task{}
	query := `SELECT * FROM tasks WHERE id=$1`
	err := s.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: task with ID %d not found", op, id)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil

}

func (s *Database) QueryTasksAll() ([]task.Task, error) {
	const op = "storage.postgresql.QueryTasksAll"

	allTasks := []task.Task{}

	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task task.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		allTasks = append(allTasks, task)
	}

	// Проверка на ошибки после завершения обработки строк.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allTasks, nil

}

func (s *Database) Close() error {
	return s.db.Close()
}

func (s *Database) UpdateTaskRecord(taskUpdate task.Task) error {
	const op = "storage.postgresql.UpdateTaskRecordByID"
	query := `UPDATE tasks SET title=$1, description=$2, due_date=$3 WHERE id=$4`
	result, err := s.db.Exec(query, taskUpdate.Title, taskUpdate.Description, taskUpdate.DueDate, taskUpdate.ID)
	fmt.Printf("Executing query: %s with ID=%d, Title=%s, Description=%s, DueDate=%s\n", query, taskUpdate.ID, taskUpdate.Title, taskUpdate.Description, taskUpdate.DueDate)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: task with ID %d not found", op, taskUpdate.ID)
	}
	return nil
}

func (s *Database) DeleteTaskRecordByID(id int) error {
	const op = "storage.postgresql.DeleteTaskRecordByID"
	query := `DELETE FROM tasks WHERE id=$1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: task with ID %d not found", op, id)
	}

	return nil
}

//DeleteTaskRecordByID(id int)
