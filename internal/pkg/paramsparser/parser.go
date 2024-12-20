package paramsparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func JSONParse[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var obj T
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&obj)
	if err != nil {
		http.Error(w, "Error decode: "+err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("error decode: %w", err)
	}
	return &obj, nil
}

func TextParse(w http.ResponseWriter, r *http.Request) (string, error) {
	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return "", err
	}

	return string(body), nil
}
