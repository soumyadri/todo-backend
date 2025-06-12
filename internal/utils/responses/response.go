package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	} else {
		w.Write([]byte("{}")) // Write an empty JSON object if no data
	}
}

func GeneralErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func ValidationErrorResponse(w http.ResponseWriter, statusCode int, errors validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errors {
		switch err.ActualTag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
			default:
				errMsgs = append(errMsgs, err.Field()+" is invalid")
		}
	}

	return Response{
		Status: "error",
		Message: "Validation failed",
		Error: strings.Join(errMsgs, ", "),
	}
}