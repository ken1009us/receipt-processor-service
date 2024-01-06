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
    retailer := receipt.Retailer
    purchaseDate := receipt.PurchaseDate
    purchaseTime := receipt.PurchaseTime
    total := receipt.Total
    items := receipt.Items

    score := 0

    // Rule 1 - One point for every alphanumeric character in the retailer name.
    for _, r := range retailer {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            score++
        }
    }

    // Rule 2 - 50 points if the total is a round dollar amount with no cents.
    totalFloat, err := strconv.ParseFloat(total, 64)
    if err == nil {
        if totalFloat == float64(int(totalFloat)) {
                score += 50
        }

        // Rule 3 - 25 points if the total is a multiple of 0.25
        if math.Mod(totalFloat, 0.25) == 0 {
            score += 25
        }
    }

    // Rule 4 - 5 points for every two items on the receipt
    numberOfItems := len(items)
    // Op / will truncate any decimal part
    score += (numberOfItems / 2) * 5

    // Rule 5 - If the trimmed length of the item description is a multiple of 3,
    // multiply the price by 0.2 and round up to the nearest integer.
    // The result is the number of points earned
    for _, item := range items {
        shortDescription := strings.TrimSpace(item.ShortDescription)
        if len(shortDescription)%3 == 0 {
            price, err := strconv.ParseFloat(item.Price, 64)
            if err == nil {
                score += int(math.Ceil(price * 0.2))
            }
        }
    }

    // Rule 6 - 6 points if the day in the purchase date is odd
    dateParts := strings.Split(purchaseDate, "-")
    if len(dateParts) == 3 {
        day, err := strconv.Atoi(dateParts[2])
        if err == nil && day % 2 != 0 {
            score += 6
        }
    }

    // Rule 7 - 10 points if the time of purchase is after 2:00pm and before 4:00pm
    timeParts, err := time.Parse("15:04", purchaseTime)
    if err == nil {
        twoPM, _ := time.Parse("15:04", "14:00")
        fourPM, _ := time.Parse("15:04", "16:00")
        if timeParts.After(twoPM) && timeParts.Before(fourPM) {
            score += 10
        }
    }

    return score
}