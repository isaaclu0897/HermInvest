package model

import (
	"encoding/json"
	"fmt"
)

// Transaction represents a record of a share transaction.
type TransactionRecord struct {
	Date      string  `gorm:"column:date"`
	Time      string  `gorm:"column:time"`
	StockNo   string  `gorm:"column:stockNo"`
	TranType  int     `gorm:"column:tranType"`
	Quantity  int     `gorm:"column:quantity"`
	UnitPrice float64 `gorm:"column:unitPrice"`
}

// NewTransactionRecord creates a new transaction record object.
func NewTransactionRecord(date, time, stockNo string, tranType, quantity int, unitPrice float64) *TransactionRecord {
	return &TransactionRecord{
		Date:      date,
		Time:      time,
		StockNo:   stockNo,
		TranType:  tranType,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func (tr *TransactionRecord) TableName() string {
	return "tblTransactionRecordSys" // default table name
}

func SumQuantityUnitPrice(remainingTrs []*TransactionRecord) (int, float64) {
	var totalQuantity, totalAmount int
	for _, tr := range remainingTrs {
		totalQuantity += tr.Quantity
		totalAmount += int(float64(tr.Quantity) * tr.UnitPrice)
	}
	avgUnitPrice := float64(totalAmount) / float64(totalQuantity)
	return totalQuantity, avgUnitPrice
}

func CalcRemainingTransactionRecords(trs []*TransactionRecord) ([]*TransactionRecord, error) {
	var remainingTrs []*TransactionRecord
	for _, tr := range trs {
		if tr.TranType > 0 {
			remainingTrs = append(remainingTrs, tr)
		} else {
			qty := tr.Quantity
			for qty > 0 && len(remainingTrs) > 0 {
				var remove *TransactionRecord
				remove, remainingTrs = remainingTrs[0], remainingTrs[1:]
				qty -= remove.Quantity
				if qty < 0 {
					return nil, fmt.Errorf("error qty < 0")

				}
			}
		}
	}
	return remainingTrs, nil
}

// Transaction represents a share transaction.
type Transaction struct {
	ID           int          `gorm:"column:id"`
	Date         string       `gorm:"column:date"`
	Time         string       `gorm:"column:time"`
	StockNo      string       `gorm:"column:stockNo"`
	TranType     int          `gorm:"column:tranType"`
	Quantity     int          `gorm:"column:quantity"`
	UnitPrice    float64      `gorm:"column:unitPrice"`
	TotalAmount  int          `gorm:"column:totalAmount"`
	Taxes        int          `gorm:"column:taxes"`
	StockMapping StockMapping `gorm:"foreignKey:stockNo;references:stockNo"`
}

// NewTransactionFromDB creates a new Transaction object from database records.
func NewTransactionFromDB(
	id int, stockNo string, date string, quantity int, tranType int,
	unitPrice float64, totalAmount int, taxes int) *Transaction {
	return &Transaction{
		ID:          id,
		StockNo:     stockNo,
		Date:        date,
		Quantity:    quantity,
		TranType:    tranType,
		UnitPrice:   unitPrice,
		TotalAmount: totalAmount,
		Taxes:       taxes,
	}
}

// NewTransactionFromInput creates a new Transaction object from input.
// It initializes the transaction with inputs. Additionally, the total amount
// and taxes are recalculated based on the new transaction details.
func NewTransactionFromInput(
	date string, time string, stockNo string, tranType int, quantity int,
	unitPrice float64) *Transaction {
	t := &Transaction{
		Date:      date,
		Time:      time,
		StockNo:   stockNo,
		TranType:  tranType,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
	t.calculateTotalAmount()
	t.calculateTaxes()
	return t
}

func (t *Transaction) TableName() string {
	return "tblTransaction" // default table name
}

// calculateTotalAmount calculates the total amount based on transaction details.
func (t *Transaction) calculateTotalAmount() {
	t.TotalAmount = int(float64(t.Quantity) * t.UnitPrice)
}

// calculateTotalAmount calculates the taxes based on transaction details.
func (t *Transaction) calculateTaxes() {
	var taxRate float64 = 0.003
	t.Taxes = int(float64(t.TotalAmount) * taxRate)
}

// SetUnitPrice updates the unit price of the transaction.
// It recalculates the total amount and taxes based on the updated unit price.
// The calculation of total amount and taxes are interdependent.
func (t *Transaction) SetUnitPrice(unitPrice float64) {
	t.UnitPrice = unitPrice

	t.recalculate()
}

// SetQuantity updates the quantity of the transaction.
// It recalculates the total amount and taxes based on the quantity.
// The calculation of total amount and taxes are interdependent.
func (t *Transaction) SetQuantity(quantity int) {
	t.Quantity = quantity

	t.recalculate()
}

// recalculate total amount and taxes of the transaction.
// It will recalculates the total amount and taxes based on the model.
func (t *Transaction) recalculate() {
	t.calculateTotalAmount()
	t.calculateTaxes()
}

type StockMapping struct {
	StockNo   string `gorm:"column:stockNo"`
	StockName string `gorm:"column:stockName"`
}

func (sp *StockMapping) TableName() string {
	return "tblStockMapping" // default table name
}

// This approach is a simple and concise way to implement custom JSON marshalling.
//
// This method is suitable for scenarios where the structure of the JSON output is
// relatively simple and the fields to be added or modified are limited and static.
func (t *Transaction) MarshalJSONbyGPT() ([]byte, error) {
	type Alias Transaction // Create an alias to avoid infinite recursion
	return json.Marshal(&struct {
		*Alias
		StockName string `json:"stock_name"`
	}{
		Alias:     (*Alias)(t),
		StockName: t.StockMapping.StockName,
	})
}

// This approach uses mapping to build a flexible struct for the JSON format.
//
// This method is more readable and provides better flexibility when you need to
// modify or add logic to the JSON output based on specific conditions.
// Reference: https://stackoverflow.com/questions/26303694/json-marshalling-unmarshalling-same-struct-to-different-json-format-in-go?rq=3
func (t *Transaction) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["ID"] = t.ID
	m["Date"] = t.Date
	m["Time"] = t.Time
	m["StockNo"] = t.StockNo
	m["StockName"] = t.StockMapping.StockName
	m["Quantity"] = t.Quantity
	m["UnitPrice"] = t.UnitPrice
	m["TotalAmount"] = t.TotalAmount
	m["Taxes"] = t.Taxes

	return json.Marshal(m)
}
