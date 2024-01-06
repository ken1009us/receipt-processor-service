package validation

import (
	"errors"
	"receipt-processor-service/model"
	"regexp"
)

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

	// Validate purchaseTime - format: "HH:MM"
	timeRegex := regexp.MustCompile("^\\d{2}:\\d{2}$")
	if !timeRegex.MatchString(receipt.PurchaseTime) {
		return errors.New("invalid purchase time format")
	}

	totalRegex := regexp.MustCompile("^\\d+\\.\\d{2}$")
	if !totalRegex.MatchString(receipt.Total) {
		return errors.New("invalid total format")
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