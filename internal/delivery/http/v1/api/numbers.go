package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

func (h *Handler) CheckNumbersForPrimeNumber(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var numbers []interface{}
	if err = json.Unmarshal(body, &numbers); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	numbersSlice := make([]uint64, len(numbers))
	for idx, number := range numbers {
		value, ok := number.(float64)
		if !ok {
			http.Error(w, fmt.Sprintf("the given input is invalid. Element on index: %d is not a number", idx),
				http.StatusBadRequest)
			return
		} else {
			// check our number has a fractional part or not
			intNumber := isFloatInt(value)
			if intNumber {
				// will be error if number is not natural number
				numbersSlice[idx] = uint64(value)
			} else {
				http.Error(w, fmt.Sprintf("the given input is invalid. Element on index: %d is not a number", idx),
					http.StatusBadRequest)
				return
			}
		}
	}

	primesSlice, err := h.numbersService.CheckNumbersForPrimeNumber(r.Context(), numbersSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(primesSlice)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonBytes))
}

func isFloatInt(floatValue float64) bool {
	return math.Mod(floatValue, 1.0) == 0
}
