package api

// This is just a small set of tests covering a few of the critical areas of functionality. In a production application, a more
// robust set of unit tests would be desired.

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetId(t *testing.T) {
	a := &Receipt{
		Retailer:      "Target",
		PurchaseDate:  "2022-01-01",
		PurchaseTime:  "13:01",
		PurchaseTotal: "35.35",
		Items:         []ReceiptItem{},
	}

	if a.Id != nil {
		t.Error("expected nil value for Id field, but received non-nil instead")
		t.Fail()
	}

	firstId := a.GetId()
	uuidValue, err := uuid.Parse(firstId)

	if err != nil {
		t.Errorf("expected GetId() method to return a valid UUID string, but encountered parsing error %s", err)
		t.Fail()
	}

	if a.Id == nil {
		t.Errorf("expected non-nil value to have been set in Id field, but received nil value instead")
		t.Fail()
	}

	if uuidValue.String() != *a.Id {
		t.Errorf("expected created UUID value to be set on Id field, but values do not match %s", *a.Id)
		t.Fail()
	}

	secondId := a.GetId()

	if firstId != secondId {
		t.Errorf("expected the same UUID to be returned by GetId(), but received another value %s (instead of %s)", secondId, firstId)
		t.Fail()
	}
}

func TestGetPurchaseDatetimeA(t *testing.T) {
	a := &Receipt{
		PurchaseDate: "2025-02-14",
		PurchaseTime: "18:14",
	}

	at, err := a.GetPurchaseDatetime()

	if err != nil {
		t.Errorf("expected nil error when parsing purchase datetime, but received error %s", err)
		t.Fail()
	}

	if at == nil || at.Year() != 2025 || at.Month() != time.February || at.Day() != 14 || at.Hour() != 18 || at.Minute() != 14 {
		t.Errorf("expected correct date time values after parsing date and time, but received incorrect values")
	}
}

func TestGetPurchaseDatetimeB(t *testing.T) {
	a := &Receipt{
		PurchaseDate: "2025-02-14",
		PurchaseTime: "18:01",
	}

	at, err := a.GetPurchaseDatetime()

	if err != nil {
		t.Errorf("expected nil error when parsing purchase datetime, but received error %s", err)
		t.Fail()
	}

	if at == nil || at.Year() != 2025 || at.Month() != time.February || at.Day() != 14 || at.Hour() != 18 || at.Minute() != 1 || at.Second() != 00 {
		t.Errorf("expected correct date time values after parsing date and time, but received incorrect values")
	}
}

//	{
//	  "retailer": "Target",
//	  "purchaseDate": "2022-01-01",
//	  "purchaseTime": "13:01",
//	  "items": [
//	    {
//	      "shortDescription": "Mountain Dew 12PK",
//	      "price": "6.49"
//	    },{
//	      "shortDescription": "Emils Cheese Pizza",
//	      "price": "12.25"
//	    },{
//	      "shortDescription": "Knorr Creamy Chicken",
//	      "price": "1.26"
//	    },{
//	      "shortDescription": "Doritos Nacho Cheese",
//	      "price": "3.35"
//	    },{
//	      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
//	      "price": "12.00"
//	    }
//	  ],
//	  "total": "35.35"
//	}
//
// Total Points: 28
// Breakdown:
//
//	   6 points - retailer name has 6 characters
//	  10 points - 5 items (2 pairs @ 5 points each)
//	   3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
//	              item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
//	   3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
//	              item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
//	   6 points - purchase day is odd
//	+ ---------
//	= 28 points
func TestGetPointsA(t *testing.T) {
	a := &Receipt{
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
	}

	points, err := a.GetPoints()

	if err != nil {
		t.Errorf("unexpected error while calaculating points. %s", err)
	}

	if points != 28 {
		t.Errorf("expected 28 points, but received %d instead.", points)
	}
}

// {
//   "retailer": "M&M Corner Market",
//   "purchaseDate": "2022-03-20",
//   "purchaseTime": "14:33",
//   "items": [
//     {
//       "shortDescription": "Gatorade",
//       "price": "2.25"
//     },{
//       "shortDescription": "Gatorade",
//       "price": "2.25"
//     },{
//       "shortDescription": "Gatorade",
//       "price": "2.25"
//     },{
//       "shortDescription": "Gatorade",
//       "price": "2.25"
//     }
//   ],
//   "total": "9.00"
// }
//
// Total Points: 109
// Breakdown:
//     50 points - total is a round dollar amount
//     25 points - total is a multiple of 0.25
//     14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
//                 note: '&' is not alphanumeric
//     10 points - 2:33pm is between 2:00pm and 4:00pm
//     10 points - 4 items (2 pairs @ 5 points each)
//   + ---------
//   = 109 points

func TestGetPointsB(t *testing.T) {
	a := &Receipt{
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
	}

	points, err := a.GetPoints()

	if err != nil {
		t.Errorf("unexpected error while calaculating points. %s", err)
	}

	if points != 109 {
		t.Errorf("expected 109 points, but received %d instead.", points)
	}
}
