package repository

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) WithTrx(trxHandle *gorm.DB) model.Repositorier {
	return &repository{db: trxHandle} // return new one
}

func (repo *repository) Begin() *gorm.DB {
	return repo.db.Begin()
}

func (repo *repository) Commit() *gorm.DB {
	return repo.db.Commit()
}

func (repo *repository) Rollback() *gorm.DB {
	return repo.db.Rollback()
}

// echo "Transaction Table" | boxes -a c -s 80 -d cc

/******************************************************************************
 *                             Transaction Table                              *
 ******************************************************************************/

// CreateTransaction: insert transaction and return inserted id
func (repo *repository) CreateTransaction(t *model.Transaction) (int, error) {
	result := repo.db.Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// CreateTransactions: insert transactions and return inserted ids
func (repo *repository) CreateTransactions(ts []*model.Transaction) ([]int, error) {
	result := repo.db.Create(&ts)
	if result.Error != nil {
		return nil, result.Error
	}

	var insertedIDs []int
	for _, t := range ts {
		insertedIDs = append(insertedIDs, t.ID)
	}

	return insertedIDs, nil
}

func (repo *repository) FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := repo.db.Where("stockNo = ?", stockNo).
		Order("date ASC, time ASC").First(&transaction).Error
	if err != nil {
		return &model.Transaction{}, err
	}

	return &transaction, nil
}

// QueryTransactionAll
func (repo *repository) QueryTransactionAll() ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := repo.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionByID
func (repo *repository) QueryTransactionByID(id int) (*model.Transaction, error) {
	var transaction *model.Transaction
	err := repo.db.Take(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// QueryTransactionByDetails
func (repo *repository) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if stockNo != "" {
		repo.db = repo.db.Where("stockNo = ?", stockNo)
	}
	if tranType != 0 {
		repo.db = repo.db.Where("tranType = ?", tranType)
	}
	if date != "" {
		repo.db = repo.db.Where("date = ?", date)
	}

	err := repo.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionInventory
func (repo *repository) QueryTransactionInventory() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	selectColumns := `stockNo, 
	sum(quantity) AS quantity, 
	sum(totalAmount)/sum(quantity) AS unitPrice, 
	sum(totalAmount) AS totalAmount, 
	sum(taxes) AS taxes`

	err := repo.db.Preload("StockMapping").
		Select(selectColumns).Group("stockNo").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionInventoryByStockNo
func (repo *repository) QueryTransactionInventoryByStockNo(stockNo string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	err := repo.db.Preload("StockMapping").
		Where("stockNo = ?", stockNo).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// updateTransaction
func (repo *repository) UpdateTransaction(id int, t *model.Transaction) error {
	err := repo.db.Updates(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) DeleteTransaction(id int) error {
	result := repo.db.Delete(&model.Transaction{ID: id})
	return result.Error
}

func (repo *repository) DeleteTransactions(ids []int) error {
	result := repo.db.Delete(&model.Transaction{}, "id IN ?", ids)
	return result.Error
}

/******************************************************************************
 *                         Transaction History Table                          *
 ******************************************************************************/

// CreateTransactionHistory: insert transaction and return inserted id
func (repo *repository) CreateTransactionHistory(t *model.Transaction) (int, error) {
	// Create a new one and set ID to 0, let SQLite autoincrement
	var transactionHistory model.Transaction = *t
	transactionHistory.ID = 0

	result := repo.db.Table("tblTransactionHistory").Create(&transactionHistory)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

/******************************************************************************
 *                                   Common                                   *
 ******************************************************************************/

// DropTable
func (repo *repository) DropTable(tablename string) error {
	// Define a whitelist of tables that are allowed to be deleted
	allowedDroppedTables := []string{
		"sqlite_sequence",
		"tblTransaction",
		"tblTransactionHistory",
		"tblTransactionCash",
		"tblTransactionRecordSys",
	}

	// Check if the tablename is in the whitelist
	found := false
	for _, allowedTable := range allowedDroppedTables {
		if allowedTable == tablename {
			found = true
			break
		}
	}

	// If the tablename is not in the whitelist, return an error
	if !found {
		return fmt.Errorf("can't drop table '%s': Table not allowed", tablename)
	}

	// Proceed with dropping the table if it's in the whitelist
	err := repo.db.Exec(fmt.Sprintf("DELETE FROM %s", tablename)).Error
	if err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                          Capital Reduction Table                           *
 ******************************************************************************/

// QueryCapitalReductionAll
func (repo *repository) QueryCapitalReductionAll() ([]*model.CapitalReduction, error) {
	var capitalReductions []*model.CapitalReduction
	if err := repo.db.Find(&capitalReductions).Error; err != nil {
		return nil, err
	}

	return capitalReductions, nil
}

// QueryCapitalReductionAll
func (repo *repository) QueryDividendAll() ([]*model.ExDividend, error) {
	var exDividends []*model.ExDividend
	if err := repo.db.Find(&exDividends).Error; err != nil {
		return nil, err
	}

	return exDividends, nil
}

/******************************************************************************
 *                          Transaction Record Table                          *
 ******************************************************************************/

// CreateTransactionRecordSys
func (repo *repository) CreateTransactionRecordSys(tr *model.TransactionRecord) error {
	if err := repo.db.Create(tr).Error; err != nil {
		return err
	}

	return nil
}

// CreateCashDividendRecord
func (repo *repository) CreateCashDividendRecord(cd *model.ExDividend) error {
	if err := repo.db.Table("tblTransactionCash").Create(cd).Error; err != nil {
		return err
	}

	return nil
}

// QueryTransactionRecordAll
func (repo *repository) QueryTransactionRecordAll() ([]*model.TransactionRecord, error) {
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Table("tblTransactionRecord").Find(&transactionRecords).Error

	if err != nil {
		return nil, nil
	}

	return transactionRecords, nil
}

// QueryTransactionRecordSysAll
func (repo *repository) QueryTransactionRecordSysAll() ([]*model.TransactionRecord, error) {
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Find(&transactionRecords).Error

	if err != nil {
		return nil, nil
	}

	return transactionRecords, nil
}

/******************************************************************************
 *                                    Note                                    *
 ******************************************************************************/

// QueryUnionNote
func (repo *repository) QueryUnionNote() {
	// Most ORMs seem not support UNION keyword, due to its complexity.
	// Faced with this situation, community suggest using "Raw" method to do this.
	// > Reference:
	// > * https://github.com/go-gorm/gorm/issues/3781
	// > * https://stackoverflow.com/questions/67190972/how-to-use-mysql-union-all-on-gorm
	// > * https://gorm.io/docs/sql_builder.html

	var transactionRecords []*model.TransactionRecord

	// Method1: Use Raw SQL with scan to query
	repo.db.Raw(`
	SELECT date, time, stockNo, tranType, quantity, unitPrice
	FROM tblTransactionRecord
	UNION SELECT * FROM tblTransactionRecordSys
	`).Scan(&transactionRecords)

	fmt.Println(transactionRecords[0], "\n\n---")

	var transactionRecords2 []*model.TransactionRecord
	// Method2: Combine GORM API build Raw SQL
	repo.db.Raw("? UNION ?",
		repo.db.Select("date, time, stockNo, tranType, quantity, unitPrice").Table("tblTransactionRecord"),
		repo.db.Select("*").Table("tblTransactionRecordSys"),
	).Scan(&transactionRecords2)

	fmt.Println(transactionRecords2[0], "\n\n---")

	// ---
	// Sometimes you would like to generate SQL without executing.
	// The Debug function is not what you expect, you can do as below.

	// Method1: Use Raw SQL without scan
	query := repo.db.Raw("SELECT * Error SQL Syntax")

	fmt.Println(query.Statement.SQL.String(), "\n\n---")

	// Method2: Use Raw SQL without scan
	stmt := repo.db.Session(&gorm.Session{DryRun: true}).Select("*").
		Table("Error SQL Syntax").Find(&model.TransactionRecord{}).Statement

	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars, "\n\n---")

	// Method3: Use ToSQL
	sql := repo.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Select("*").Table("Error SQL Syntax").Find(&model.TransactionRecord{})
	})

	fmt.Println(sql)
}

// QueryTransactionRecordSys
func (repo *repository) QueryTransactionPreload() ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := repo.db.Preload("StockMapping").Take(&transactions).Error
	if err != nil {
		return nil, nil
	}

	return transactions, nil
}
