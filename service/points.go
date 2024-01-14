package service

import (
	"math"
	"receipt-processor-service/model"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func CalculatePoints(receipt *model.Receipt) int {
    score := 0

    score += calculatePointsForRetailerName(receipt.Retailer)
    score += calculatePointsForTotalAmount(receipt.Total)
    score += calculatePointsPerItem(receipt.Items)
    score += calculatePointsForPurchaseDate(receipt.PurchaseDate)
    score += calculatePointsForPurchaseTime(receipt.PurchaseTime)

    return score
}

func calculatePointsForRetailerName(retailer string) int {
    points := 0
    for _, r := range retailer {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            points++
        }
    }
    return points
}

func calculatePointsForTotalAmount(total string) int {
    points := 0
    totalFloat, err := strconv.ParseFloat(total, 64)
    if err == nil {
        if totalFloat == float64(int(totalFloat)) {
            points += 50
        }
        if math.Mod(totalFloat, 0.25) == 0 {
            points += 25
        }
    }
    return points
}

func calculatePointsPerItem(items []model.Item) int {
    points := len(items) / 2 * 5
    for _, item := range items {
        points += calculatePointsForItemDescription(item.ShortDescription, item.Price)
    }
    return points
}

func calculatePointsForItemDescription(description, price string) int {
    if len(strings.TrimSpace(description))%3 == 0 {
        priceFloat, err := strconv.ParseFloat(price, 64)
        if err == nil {
            return int(math.Ceil(priceFloat * 0.2))
        }
    }
    return 0
}

func calculatePointsForPurchaseDate(purchaseDate string) int {
    dateParts := strings.Split(purchaseDate, "-")
    if len(dateParts) == 3 {
        day, err := strconv.Atoi(dateParts[2])
        if err == nil && day%2 != 0 {
            return 6
        }
    }
    return 0
}

func calculatePointsForPurchaseTime(purchaseTime string) int {
    timeParts, err := time.Parse("15:04", purchaseTime)
    if err == nil {
        twoPM, _ := time.Parse("15:04", "14:00")
        fourPM, _ := time.Parse("15:04", "16:00")
        if timeParts.After(twoPM) && timeParts.Before(fourPM) {
            return 10
        }
    }
    return 0
}
