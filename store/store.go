package store

import (
	"fmt"
	"receipt-processor-service/model"
	"sync"
)


var receiptsMap = make(map[string]*model.Receipt)
var mapMutex = &sync.Mutex{}

func StoreReceipt(id string, receipt *model.Receipt) {
	mapMutex.Lock()
	receiptsMap[id] = receipt
	mapMutex.Unlock()
}

func RetrieveReceipt(id string) (*model.Receipt, bool) {
	mapMutex.Lock()
	receipt, exists := receiptsMap[id]
	mapMutex.Unlock()
	fmt.Printf("Attempting to retrieve receipt with ID: %s\n", id)
	fmt.Printf("Found: %v\n", exists)
	return receipt, exists
}