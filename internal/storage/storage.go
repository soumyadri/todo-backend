package storage

import "github.com/soumyadri/todo-backend/internal/types"

type Storage interface {
	NewTodos(todo types.Todos) (int64, error)
	GetAllTodos() ([]types.Todos, error)
	GetTodoByStatus(status string) ([]types.Todos, error)
	GetTodoByDoneBy() ([]types.Todos, error)
	GetTodoById(id int64) (types.Todos, error)
	UpdateTodo(id int64, todo types.Todos) error
}