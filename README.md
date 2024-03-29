# Receipt Processor Service

This repository contains the Receipt Processor Service, a Go-based web service for processing receipts and calculating points based on various criteria.

## Prerequisites

- Go (1.20)
- Docker

## Installation

1. Clone this repository:

```bash
$ git clone https://github.com/ken1009us/receipt-processor-service.git
```

2. Navigate to the project directory:

```bash
$ cd receipt-processor-service
```

3. Installing Go

Follow the instructions at https://golang.org/doc/install to download and install Go.

## Usage

### Running Locally

1. Navigate to the project directory.

2. To start the service, run:

``` bash
$ go run .
```

This will start the server on http://localhost:8080.

### Using Docker (Recommended)

1. Make sure you have Docker installed on your system. You can download and install Docker from the official Docker website: https://www.docker.com/.

2. Navigate to the project directory in your terminal.

3. Build the Docker image using the following command:

```bash
$ docker build -t receipt-processor-service .
```

This command builds a Docker image with the name receipt-processor based on the Dockerfile in the current directory. The -t flag specifies the image name.

4. After the image is built successfully, you can run a Docker container using the image with the following command:

```bash
$ docker run -d --name myreceiptprocessor -p 8080:8080 receipt-processor-service
```

This command creates and starts a new Docker container named myreceiptprocessor from the receipt-processor-service image. It runs the container in detached mode (in the background), and maps the container's port 8080 to the host's port 8080.

## Process the receipt

1. To process a receipt, use the following curl command:

```bash
$ curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d
'{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
```

This will return an ID for the processed receipt.

```bash
{"id":"4a399209-8f6b-4ff6-ab72-0dbcfe8efb50"}
```

2. To retrieve the points for the receipt, use the ID returned from the above command in the following curl command:

Replace {RECEIPT_ID} with the actual ID.

```bash
$ curl http://localhost:8080/receipts/{RECEIPT_ID}/points
```

This should return the calculated points in JSON format, such as:

```bash
{"points":28}
```

## Project Structure

- main.go: Entry point for the web service.
- handler/: Contains HTTP request handlers.
- model/: Data models for receipts and items.
- service/: Business logic for calculating points.
- store/: In-memory storage and logic for receipts.
- util/: Error handling and other utilities functions.
- validation/: Validates the format of input files.
- Dockerfile: Instructions for building the Docker image.

## Improvements

There is another way to handle the calculated points.

### Extending the Receipt Struct:

Advantages:
- Cohesion: Integrates points directly with receipts, making the data model more intuitive and easier to understand.
- Simplicity in Data Retrieval: Simplifies data access patterns, as points are directly associated with their respective receipts.
- Database Storage: Easier to manage in a database, as each receipt row can include the points without requiring joins or additional lookups.

Disadvantages:
- Increased Data Payload: If there are scenarios where you need to access receipt data without points, you'll still be carrying the points data, which could be inefficient.
- Potential for Redundant Calculations: If points calculations are complex and you store the calculated points with each receipt, you might end up recalculating points unnecessarily in some scenarios.

In the model.go file:

```go
package model

type Receipt struct {
    // Existing fields...
    Retailer       string `json:"retailer"`
    PurchaseDate   string `json:"purchaseDate"`
    PurchaseTime   string `json:"purchaseTime"`
    Total          string `json:"total"`
    Items          []Item `json:"items"`
    // Add a new field for points
    Points         int    `json:"points"`
}

```

In the handler file:

```go
func ProcessReceipt(receipt *model.Receipt) {
    // ...
    // Calculate points
    receipt.Points = service.CalculatePoints(receipt)

    // Store the receipt with points
    // ...
}
```

Then we can use API endpoint to retrieve receipts.

```go
func GetReceipt(id string) (*model.Receipt, error) {
    // ...
    receipt, exists := store.RetrieveReceipt(id)
    if !exists {
        return nil, errors.New("receipt not found")
    }
    return receipt, nil
}
```