package validation

import (
	"errors"
	"receipt-processor-service/model"
	"regexp"
	"strconv"
	"time"
)

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func ValidateReceiptTotal(receipt model.Receipt) error {
	var sum float64 = 0
	for _, item := range receipt.Items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return errors.New("invalid price format")
		}

		sum += price
	}

	receiptTotal, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        return errors.New("invalid total format")
    }

	const t = 0.01
	if abs(receiptTotal - sum) > t {
        return errors.New("the sum of item prices does not match total")
    }

	return nil
}

func ValidateReceipt(receipt model.Receipt) error {
	// Validate retailer - non-empty, no leading/trailing white spaces
	if match, _ := regexp.MatchString("^[\\w&\\s-]+$", receipt.Retailer); !match {
		return errors.New("invalid retailer format")
	}

	// Validate purchaseDate - format: "YYYY-MM-DD"
	dateRegex := regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$")
	if !dateRegex.MatchString(receipt.PurchaseDate) {
		return errors.New("invalid purchase date format")
	}

	// Validate purchaseDate for logical correctness
	_, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return errors.New("invalid purchase date: " + err.Error())
	}

	// Validate purchaseTime - format: "HH:MM"
	timeRegex := regexp.MustCompile("^\\d{2}:\\d{2}$")
	if !timeRegex.MatchString(receipt.PurchaseTime) {
		return errors.New("invalid purchase time format")
	}

	// Validate purchaseTime for logical correctness
	_, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return errors.New("invalid purchase time: " + err.Error())
	}

	totalRegex := regexp.MustCompile("^\\d+\\.\\d{2}$")
	if !totalRegex.MatchString(receipt.Total) {
		return errors.New("invalid total format")
	}

	err = ValidateReceiptTotal(receipt)
	if err != nil {
		return errors.New("invalid receipt total and the sum of the items price: " + err.Error())
	}

	for _, item := range receipt.Items {
		desRegex := regexp.MustCompile("^[\\w\\s\\-]+$")
		if !desRegex.MatchString(item.ShortDescription) {
			return errors.New("invalid item short description format")
		}

		if !totalRegex.MatchString(item.Price) {
			return errors.New("invalid item price format")
		}
	}

	return nil
}