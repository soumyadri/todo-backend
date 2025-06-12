package storage

import "github.com/soumyadri/todo-backend/internal/types"

type Storage interface {
	NewTodos(todo types.Todos) (int64, error)
	GetAllTodos() ([]types.Todos, error)
}