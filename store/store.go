package store

import (
	"fmt"
	"receipt-processor-service/model"
	"sync"
)


var receiptsMap = make(map[string]*model.Receipt)
var pointsMap = make(map[string]int)
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

func StorePoint(id string, point int) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	pointsMap[id] = point
}

func RetrievePoint(id string) (int, bool) {
	mapMutex.Lock()
	point, exists := pointsMap[id]
	mapMutex.Unlock()

	fmt.Printf("Attempting to retrieve point with ID: %s\n", id)
	fmt.Printf("Found: %v\n", exists)

	return point, exists

}