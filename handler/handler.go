package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-processor-service/model"
	"receipt-processor-service/service"
	"receipt-processor-service/store"
	"receipt-processor-service/validation"

	"strings"

	"github.com/google/uuid"
)

// ProcessReceipt  - POST /receipts/process
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------Process Receipt------")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Bad request: invalid receipt format", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := validation.ValidateReceipt(receipt); err != nil {
		http.Error(w, "Invalid receipt: "+err.Error(), http.StatusBadRequest)
        return
	}

	id := uuid.New().String()
	store.StoreReceipt(id, &receipt)
	fmt.Printf("Stored receipt with ID: %s\n", id)
    fmt.Printf("Receipt: %+v\n", receipt)

    json.NewEncoder(w).Encode(model.ReceiptResponse{ID: id})
}


// GetPoints - GET /receipts/{id}/points
func GetPoints(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------Retrieve Receipt------")
	pathSegments := strings.Split(r.URL.Path, "/")
    if len(pathSegments) < 3 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    id := pathSegments[2]

	fmt.Printf("Retrieving receipt with ID: %s\n", id)

	receipt, exists := store.RetrieveReceipt(id)

	if !exists {
		http.Error(w, "Not found: receipt ID does not exist", http.StatusNotFound)
		return
	}

	point := service.CalculatePoints(receipt)
	json.NewEncoder(w).Encode(model.PointsResponse{Points: point})

}

