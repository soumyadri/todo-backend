package todo

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/soumyadri/todo-backend/internal/storage"
	"github.com/soumyadri/todo-backend/internal/types"
	"github.com/soumyadri/todo-backend/internal/utils/responses"
)

func NewTodos(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request for creating new todos", "method", r.Method, "path", r.URL.Path)

		var todo types.Todos
		err := json.NewDecoder(r.Body).Decode(&todo)
		if errors.Is(err, io.EOF) {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "No todo data provided")
			return
		}
		if err != nil {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Invalid todo details data")
			return
		}

		defer r.Body.Close()

		// Validate the todo data
		error := validator.New().Struct(&todo)
		if error != nil {
			slog.Error("Validation failed for todo data", "error", error)
			if validationErrors, ok := error.(validator.ValidationErrors); ok {
				// Handle validation errors
				errResponse := response.ValidationErrorResponse(w, http.StatusBadRequest, validationErrors)
				slog.Error("Validation errors", "errors", errResponse)
				response.WriteJson(w, http.StatusBadRequest, errResponse)
			} else {
				// Handle other types of errors
				response.GeneralErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				slog.Error("Unexpected error during validation", "error", error)
				return
			}
			return
		}

		todoId, error := storage.NewTodos(todo)
		if error != nil {
			slog.Error("Failed to create todo in storage", "error", error)
			response.GeneralErrorResponse(w, http.StatusInternalServerError, "Failed to create todo")
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Todo data received successfully", "id": strconv.FormatInt(todoId, 10) })
		slog.Info("Response sent successfully", "todo", todoId)
	}
}

func GetTodos(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Get All the Todos", "method", r.Method, "path", r.URL.Path, "query", r.URL.Query())
		
		var todos []types.Todos
		var err error

		if r.URL.Query().Get("status") != "" {
			status := r.URL.Query().Get("status")
			todos, err = storage.GetTodoByStatus(status)
		} else if r.URL.Query().Get("doneby") == "upcoming" {
			todos, err = storage.GetTodoByDoneBy()
		} else {
			todos, err = storage.GetAllTodos()
		}
		
		if err != nil {
			slog.Error("Failed to retrieve todos from storage", "error", err)
			response.GeneralErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos")
			return
		}

		if len(todos) == 0 {
			slog.Info("No todos found")
			response.WriteJson(w, http.StatusOK, map[string]string{"status": "success", "message": "No todos found"})
			return
		}

		response.WriteJson(w, http.StatusOK, todos)
		slog.Info("Successfully retrieved students", "count", len(todos))
	}
}

func GetTodoById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Get Todo by ID", "method", r.Method, "path", r.URL.Path)

		idStr := r.PathValue("id")
		if idStr == "" {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Todo ID is required")
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error("Invalid Todo ID format", "error", err)
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Invalid Todo ID format")
			return
		}

		todo, err := storage.GetTodoById(id)
		if err != nil {
			slog.Error("Failed to retrieve todo from storage", "error", err)
			response.GeneralErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo")
			return
		}

		if (todo == types.Todos{}) {
			slog.Info("Todo not found", "id", id)
			response.GeneralErrorResponse(w, http.StatusNotFound, "Todo not found")
			return
		}

		response.WriteJson(w, http.StatusOK, todo)
		slog.Info("Successfully retrieved todo by ID", "id", id)
	}
}

func UpdateTodo(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Update Todo by ID", "method", r.Method, "path", r.URL.Path)

		idStr := r.PathValue("id")
		if idStr == "" {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Todo ID is required")
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error("Invalid Todo ID format", "error", err)
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Invalid Todo ID format")
			return
		}

		var todo types.Todos
		err = json.NewDecoder(r.Body).Decode(&todo)
		if errors.Is(err, io.EOF) {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "No todo data provided")
			return
		}
		if err != nil {
			response.GeneralErrorResponse(w, http.StatusBadRequest, "Invalid todo details data")
			return
		}

		defer r.Body.Close()

		todo.ID = int(id) // Set the ID for the todo to be updated

		error := validator.New().Struct(&todo)
		if error != nil {
			slog.Error("Validation failed for todo data", "error", error)
			if validationErrors, ok := error.(validator.ValidationErrors); ok {
				errResponse := response.ValidationErrorResponse(w, http.StatusBadRequest, validationErrors)
				slog.Error("Validation errors", "errors", errResponse)
				response.WriteJson(w, http.StatusBadRequest, errResponse)
			} else {
				response.GeneralErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				slog.Error("Unexpected error during validation", "error", error)
				return
			}
			return
		}

		err = storage.UpdateTodo(id, todo)
		if err != nil {
			slog.Error("Failed to update todo in storage", "error", err)
			response.GeneralErrorResponse(w, http.StatusInternalServerError, "Failed to update todo")
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Todo updated successfully"})
		slog.Info("Successfully updated todo", "id", id)
	}
}