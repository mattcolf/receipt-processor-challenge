package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-memdb"
)

type ReceiptDatabase struct {
	MemDB *memdb.MemDB
}

// Setup and initialize the MemDB database
func SetupDatabase() *ReceiptDatabase {
	// Setup a basic schema that enables querying of receipts by Id
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"receipt": {
				Name: "receipt",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	memdb, err := memdb.NewMemDB(schema)

	if err != nil {
		log.Fatalf("Error while initializing database. %s", err)
	}

	db := &ReceiptDatabase{
		MemDB: memdb,
	}

	// Load some sample data to make it easier to test
	db.LoadExampleData()

	return db
}

// Insert a new receipt
func (db ReceiptDatabase) InsertReceipt(receipt *Receipt) (*string, error) {
	// ensure that the ID is set before database insertion
	id := receipt.GetId()

	if receipt.Id == nil {
		return nil, errors.New("unable to insert receipt with no Id field value set")
	}

	txn := db.MemDB.Txn(true)

	if err := txn.Insert("receipt", receipt); err != nil {
		txn.Abort()

		return nil, fmt.Errorf("unable to insert receipt because of unknown error. %s", err)
	}

	txn.Commit()

	return &id, nil
}

// Get all receipts
func (db ReceiptDatabase) GetAllReceipts() ([]*Receipt, error) {
	receipts := make([]*Receipt, 0)

	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("receipt", "id")

	if err != nil {
		return nil, fmt.Errorf("error while querying receipts from database. %s", err)
	}

	// TODO: pagination, offsets, filtering, etc.
	for raw := it.Next(); raw != nil; raw = it.Next() {
		receipt := raw.(*Receipt)
		receipts = append(receipts, receipt)
		log.Printf("loaded receipt with id %s", receipt.GetId())
	}

	return receipts, nil
}

// Get a receipt by ID
func (db ReceiptDatabase) GetReceiptById(id string) (*Receipt, error) {
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("receipt", "id", id)

	if err != nil {
		return nil, fmt.Errorf("error while querying database for id %s. %s", id, err)
	}

	return raw.(*Receipt), nil
}

// Load example data into the database
func (db ReceiptDatabase) LoadExampleData() {
	txn := db.MemDB.Txn(true)

	// override the default id generation to ensure consistent Id values across server restarts for example data
	idA := "ae6ad71c-e978-4e93-a8a9-5909dc3d4422"
	idB := "98557e85-1663-4d3b-adff-bbba1d002c4e"
	idC := "392abbcf-4783-49f4-901c-ae0c708783df"
	idD := "cbf19128-6408-4b47-9d20-08c2e84a9341"

	receipts := []*Receipt{
		{
			Id:            &idA,
			Retailer:      "Walgreens",
			PurchaseDate:  "2022-01-02",
			PurchaseTime:  "08:13",
			PurchaseTotal: "2.65",
			Items: []ReceiptItem{
				{
					ShortDescription: "Pepsi - 12-oz",
					Price:            "1.25",
				},
				{
					ShortDescription: "Dasani",
					Price:            "1.40",
				},
			},
		},
		{
			Id:            &idB,
			Retailer:      "Target",
			PurchaseDate:  "2022-01-02",
			PurchaseTime:  "13:13",
			PurchaseTotal: "1.25",
			Items: []ReceiptItem{
				{
					ShortDescription: "Pepsi - 12-oz",
					Price:            "1.25",
				},
			},
		},
		{
			Id:            &idC,
			Retailer:      "Target",
			PurchaseDate:  "2022-01-01",
			PurchaseTime:  "13:01",
			PurchaseTotal: "35.35",
			Items: []ReceiptItem{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
				{
					ShortDescription: "Knorr Creamy Chicken",
					Price:            "1.26",
				},
				{
					ShortDescription: "Doritos Nacho Cheese",
					Price:            "3.35",
				},
				{
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            "12.00",
				},
			},
		},
		{
			Id:            &idD,
			Retailer:      "M&M Corner Market",
			PurchaseDate:  "2022-03-20",
			PurchaseTime:  "14:33",
			PurchaseTotal: "9.00",
			Items: []ReceiptItem{
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
			},
		},
	}

	for _, receipt := range receipts {
		id := receipt.GetId()
		log.Printf("inserting example receipt with id %s", id)
		if err := txn.Insert("receipt", receipt); err != nil {
			log.Fatalf("Error while inserting example database data. %s", err)
		}
	}

	txn.Commit()
}
