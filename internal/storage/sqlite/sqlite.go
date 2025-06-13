package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/soumyadri/todo-backend/internal/config"
	"github.com/soumyadri/todo-backend/internal/types"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL,
		duedate DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	
	// Check for errors in table creation
	if err != nil {
		log.Println("Error creating students table:", err)
		return nil, err
	}

	return &SQLite{Db: db}, nil
}

func (db *SQLite) NewTodos(todo types.Todos) (int64, error) {
	result, err := db.Db.Exec("INSERT INTO todos (title, description, status, dueDate) VALUES (?,?,?,?)", todo.Title, todo.Description, todo.Status, todo.Duedate);
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId();
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *SQLite) GetAllTodos() ([]types.Todos, error) {
	rows, err := db.Db.Query("SELECT id, title, description, status, duedate, created_at, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []types.Todos
	for rows.Next() {
		var todo types.Todos
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Duedate, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (db *SQLite) GetTodoByStatus(status string) ([]types.Todos, error) {
	rows, err := db.Db.Query("SELECT id, title, description, status, duedate, created_at, updated_at FROM todos WHERE status = ?", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []types.Todos
	for rows.Next() {
		var todo types.Todos
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Duedate, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (db *SQLite) GetTodoByDoneBy() ([]types.Todos, error) {
	rows, err := db.Db.Query("SELECT * FROM todos WHERE dueDate > ? ORDER BY dueDate ASC;", time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []types.Todos
	for rows.Next() {
		var todo types.Todos
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Duedate, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (db *SQLite) GetTodoById(id int64) (types.Todos, error) {
	row := db.Db.QueryRow("SELECT id, title, description, status, duedate, created_at, updated_at FROM todos WHERE id = ?", id)

	var todo types.Todos
	if err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Duedate, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return todo, nil // No todo found with the given ID
		}
		return todo, err // Other error
	}

	return todo, nil
}

func (db *SQLite) UpdateTodo(id int64, todo types.Todos) error {
	_, err := db.Db.Exec("UPDATE todos SET title = ?, description = ?, status = ?, duedate = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", todo.Title, todo.Description, todo.Status, todo.Duedate, id)
	if err != nil {
		return err
	}

	return nil
}