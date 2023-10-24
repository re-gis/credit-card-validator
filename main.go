package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/validate", ValidateFunc)
	http.ListenAndServe(":8080", nil)
}

func IsValidLuhn(s string) bool {
	var sum int
	var alternate bool
	for i := len(s) - 1; i >= 0; i-- {
		n := int(s[i] - '0')
		if alternate {
			n *= 2
			if n > 9 {
				n = (n % 10) + 1
			}
		}
		sum += n
		alternate = !alternate
	}
	return sum%10 == 0
}

func ValidateFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method are allowed!", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Number string `json:"number"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}

	if payload.Number == "" {
		http.Error(w, "Missing 'number' in payload", http.StatusBadRequest)
		return
	}

	isValid := IsValidLuhn(payload.Number)

	// Wrapping the result into a JSON response payload
	result := struct {
		IsValid bool `json:"isValid"`
	}{
		IsValid: isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
