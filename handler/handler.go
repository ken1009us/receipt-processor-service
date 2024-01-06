package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"receipt-processor-service/model"
	"receipt-processor-service/service"
	"receipt-processor-service/store"
	"receipt-processor-service/util"
	"receipt-processor-service/validation"

	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// ProcessReceipt  - POST /receipts/process
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Processing receipt")

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		util.WriteError(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		util.WriteError(w, "Bad request: invalid receipt format", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := validation.ValidateReceipt(receipt); err != nil {
		util.WriteError(w, "Invalid receipt: " + err.Error(), http.StatusBadRequest)
        return
	}

	id := uuid.New().String()
	store.StoreReceipt(id, &receipt)

    json.NewEncoder(w).Encode(model.ReceiptResponse{ID: id})
}


// GetPoints - GET /receipts/{id}/points
func GetPoints(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Retrieving receipt")

	w.Header().Set("Content-Type", "application/json")

	pathSegments := strings.Split(r.URL.Path, "/")
    if len(pathSegments) < 3 {
		util.WriteError(w, "Invalid path", http.StatusBadRequest)
        return
    }
    id := pathSegments[2]

	receipt, exists := store.RetrieveReceipt(id)
	if !exists {
		util.WriteError(w, "Not found: receipt ID does not exist!", http.StatusNotFound)
		return
	}

	point := service.CalculatePoints(receipt)
	json.NewEncoder(w).Encode(model.PointsResponse{Points: point})

}

