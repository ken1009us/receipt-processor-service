package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt-processor-service/handler"
)

func main() {
	http.HandleFunc("/receipts/process", handler.ProcessReceipt)
	http.HandleFunc("/receipts/", handler.GetPoints)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
