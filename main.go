package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prac/calc"
)

func calculateHandler(calculator *calc.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var instructions []calc.Instruction
		if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		results, err := calculator.Calculate(instructions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"items": results,
		})
	}
}

func main() {
	calculator := calc.NewCalculator()

	http.HandleFunc("/calculate", calculateHandler(calculator))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
