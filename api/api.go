package api

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReceiptsApi struct {
	Router   *gin.Engine
	Database *ReceiptDatabase
}

func SetupApi() *ReceiptsApi {
	api := &ReceiptsApi{
		Router:   gin.Default(),
		Database: SetupDatabase(),
	}

	api.Router.NoRoute(api.HandleNoRoute)
	api.Router.NoMethod(api.HandleNoMethod)
	api.Router.GET("/status", api.HandleStatus)

	api.Router.POST("/receipts/process", api.HandleCreateNewReceipt)
	api.Router.GET("/receipts", api.HandleGetAllReceipts)
	api.Router.GET("/receipts/:id", api.HandleGetReceiptById)
	api.Router.GET("/receipts/:id/points", api.HandleGetReceiptPointsById)

	return api
}

type Header struct {
	ContentType *string `header:"content-type" binding:"required"`
}

// Create a new receipt record
// POST /receipts/process
func (api ReceiptsApi) HandleCreateNewReceipt(c *gin.Context) {

	var input Receipt
	var header Header

	// bind to headers and perform basic check
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(400, gin.H{
			"error": "Missing content-type header.",
		})

		return
	}

	if header.ContentType == nil || !strings.Contains(*header.ContentType, "application/json") {
		c.JSON(405, gin.H{
			"error": "Unsupported content type. Only `application/json` is supported.",
		})

		return
	}

	// bind to input and perform basic format validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "The receipt is invalid.",
		})

		return
	}

	errors := make([]string, 0)

	// validate retailer
	if match := regexp.MustCompile(`^[\w\s\-&]+$`).MatchString(input.Retailer); !match {
		errors = append(errors, "invalid retailer")
	}

	// validate purchase date
	if match := regexp.MustCompile(`^[0-9]{4}\-[0-1][0-9]\-[0-3][0-9]$`).MatchString(input.PurchaseDate); !match {
		errors = append(errors, "invalid purchaseDate value")
	}

	// validate purchase time
	if match := regexp.MustCompile(`^[0-2][0-9]:[0-5][0-9]$`).MatchString(input.PurchaseTime); !match {
		errors = append(errors, "invalid purchaseTime value")
	}

	// validate total
	if match := regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(input.PurchaseTotal); !match {
		errors = append(errors, "invalid total value")
	}

	// validate receipt item count
	if len(input.Items) < 1 {
		errors = append(errors, "at least one receipt item must be provided")
	}

	// validate each receipt item
	for i, item := range input.Items {
		// validate receipt item description
		if match := regexp.MustCompile(`^[\w\s\-]+$`).MatchString(item.ShortDescription); !match {
			errors = append(errors, fmt.Sprintf("invalid shortDescription value for receipt item %d", i))
		}

		// validate receipt item price
		if match := regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price); !match {
			errors = append(errors, fmt.Sprintf("invalid price value for receipt item %d", i))
		}
	}

	// return all validation errors if any were encountered
	if len(errors) > 0 {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("The receipt is invalid. %s", strings.Join(errors, ", ")),
		})

		return
	}

	// insert a new receipt record to the database
	id, err := api.Database.InsertReceipt(&input)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "unknown error",
		})

		log.Printf("error while inserting receipt record into database. %s", err)
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

// Query all receipts
// GET /receipts
func (api ReceiptsApi) HandleGetAllReceipts(c *gin.Context) {
	receipts, err := api.Database.GetAllReceipts()

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while querying all receipts",
		})

		log.Printf("error while querying all receipts. %s", err)
		return
	}

	c.JSON(200, receipts)
}

// Query a single receipt by ID
// GET /receipts/{id}
func (api ReceiptsApi) HandleGetReceiptById(c *gin.Context) {
	id := c.Param("id")

	// validate the id format (GUID)
	if match := regexp.MustCompile(`^[{]?[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}[}]?$`).MatchString(id); !match {
		c.JSON(400, gin.H{
			"error": "invalid receipt id format",
		})

		log.Printf("invalid id %s provided", id)
		return
	}

	// lookup receipt by ID
	receipt, err := api.Database.GetReceiptById(id)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "no receipt found",
		})

		log.Printf("no receipt found for id %s. %s", id, err)
		return
	}

	c.JSON(200, receipt)
}

// Query a single receipt by ID and return the number of points
// GET /receipts/{id}/points
func (api ReceiptsApi) HandleGetReceiptPointsById(c *gin.Context) {
	id := c.Param("id")

	// validate the id format (GUID)
	if match := regexp.MustCompile(`^[{]?[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}[}]?$`).MatchString(id); !match {
		c.JSON(400, gin.H{
			"error": "invalid receipt id format",
		})

		log.Printf("invalid id %s provided", id)
		return
	}

	// lookup receipt by ID
	receipt, err := api.Database.GetReceiptById(id)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "no receipt found",
		})

		log.Printf("no receipt found for id %s. %s", id, err)
		return
	}

	// calculate receipt points
	points, err := receipt.GetPoints()

	if err != nil {
		c.JSON(500, gin.H{
			"error": "unknown error while calculating receipt points",
		})

		log.Printf("error while calculating points for id %s. %s", id, err)
		return
	}

	c.JSON(200, gin.H{
		"points": points,
	})
}

func (api ReceiptsApi) HandleNoRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"error": "Route not found.",
	})
}

func (api ReceiptsApi) HandleNoMethod(c *gin.Context) {
	c.JSON(405, gin.H{
		"error": "Method not allowed.",
	})
}

func (api ReceiptsApi) HandleStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
