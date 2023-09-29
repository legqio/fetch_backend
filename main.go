package main

import (
	"errors"
	"legqio/backend_challenge/pkg/receipts"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()

	// Instantiate receipt Handler and provide a data store
	store := receipts.NewMemStore()
	receiptsHandler := NewReceiptsHandler(store)

	// Register Routes
	router.POST("/receipts/process", receiptsHandler.processReceipts)
	router.GET("/receipts/:id/points", receiptsHandler.getPoints)

	router.Run()
}

type ReceiptsHandler struct {
	store receiptStore
}

func NewReceiptsHandler(s receiptStore) *ReceiptsHandler {
	return &ReceiptsHandler{
		store: s,
	}
}

type receiptStore interface {
	PostReceipt(name string, receipt receipts.Receipt, points int) error
	GetPoints(name string) (int, error)
}

func validateReceipt(receipt receipts.Receipt) error {

	if len(receipt.Retailer) == 0 {
		return errors.New("receipt must include a retailer")
	}

	if len(receipt.PurchaseDate) == 0 {
		return errors.New("receipt must include a purchase date")
	}

	if len(receipt.PurchaseTime) == 0 {
		return errors.New("your receipt must include a purchase time")
	}

	if len(receipt.Total) == 0 {
		return errors.New("your receipt must include a total cost")
	}

	if len(receipt.Items) == 0 {
		return errors.New("your receipt must include at least one item")
	}

	for _, item := range receipt.Items {
		if len(item.ShortDescription) == 0 {
			return errors.New("all receipt items must include a description")
		}
		if len(item.Price) == 0 {
			return errors.New("all receipt items must include a price")
		}
		_, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return err
		}
	}

	return nil
}

func calcPoints(receipt receipts.Receipt) (points int, err error) {
	points = 0

	regex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(regex.FindAllString(receipt.Retailer, -1))

	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0, err
	}

	if math.Mod(totalFloat, 1) == 0 {
		points += 50
	}

	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if math.Mod(float64(len(strings.TrimSpace(item.ShortDescription))), 3) == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}

			points += int(math.Ceil(priceFloat * 0.2))
		}
	}

	layoutDate := "2006-01-02"
	parsedPurchaseDate, err := time.Parse(layoutDate, receipt.PurchaseDate)
	if err != nil {
		return 0, err
	}
	if math.Mod(float64(parsedPurchaseDate.Day()), 2) == 1 {
		points += 6
	}

	layoutTime := "15:04"
	parsedPurchaseTime, err := time.Parse(layoutTime, receipt.PurchaseTime)
	if err != nil {
		return 0, err
	}
	minTime, err := time.Parse(layoutTime, "14:00")
	if err != nil {
		return 0, err
	}
	maxTime, err := time.Parse(layoutTime, "16:00")
	if err != nil {
		return 0, err
	}
	if parsedPurchaseTime.Sub(minTime).Minutes() > 0 && maxTime.Sub(parsedPurchaseTime).Minutes() > 0 {
		points += 10
	}

	return points, nil
}

func (h ReceiptsHandler) processReceipts(c *gin.Context) {
	var receipt receipts.Receipt
	err := c.ShouldBindJSON(&receipt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate incoming receipt properties and structure
	errValidateReceipt := validateReceipt(receipt)
	if errValidateReceipt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errValidateReceipt.Error()})
		return
	}

	// Calculate points
	points, errCalcPoints := calcPoints(receipt)
	if errCalcPoints != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCalcPoints.Error()})
		return
	}

	// Create receipt id (UUID)
	id := uuid.New()

	// Save receipt in memory storage
	h.store.PostReceipt(id.String(), receipt, points)

	c.JSON(
		http.StatusOK,
		gin.H{
			"id": id.String(),
		},
	)
}

func (h ReceiptsHandler) getPoints(c *gin.Context) {
	id := c.Param("id")

	points, err := h.store.GetPoints(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"points": points,
			},
		)
	}
}
