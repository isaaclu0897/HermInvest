package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type transactionRepository struct {
	db *sql.DB
}

// CreateTransaction: insert transaction and return inserted id
func (repo *transactionRepository) createTransaction(t *Transaction) (int, error) {
	query := `INSERT INTO tblTransaction (stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := repo.db.Exec(query, t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
	if err != nil {
		fmt.Println("Error insert database: ", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting inserted id: ", err)
		return 0, err
	}
	return int(id), nil
}

// testcase begin, commit, rollback
// CreateTransactions: insert transaction and return inserted id
func (repo *transactionRepository) createTransactions(ts []*Transaction) error {

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO tblTransaction (stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range ts {
		rst, err := stmt.Exec(t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
		if err != nil {
			return err
		}
		_, err = rst.LastInsertId()
		if err != nil {
			fmt.Println("Error getting inserted id: ", err)
			return err
		}
		// fmt.Printf("inserted id: %d\n", id)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// queryTransactionAll
func (repo *transactionRepository) queryTransactionAll() ([]*Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// queryTransactionByID
func (repo *transactionRepository) queryTransactionByID(id int) ([]*Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var transactions []*Transaction
	var t Transaction
	err := row.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, &t)

	return transactions, nil
}

// queryTransactionByDetails
func (repo *transactionRepository) queryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error) {
	var conditions []string
	var args []interface{}

	if stockNo != "" {
		conditions = append(conditions, "stockNo = ?")
		args = append(args, stockNo)
	}
	if tranType != 0 {
		conditions = append(conditions, "tranType = ?")
		args = append(args, tranType)
	}
	if date != "" {
		conditions = append(conditions, "date = ?")
		args = append(args, date)
	}

	query := fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// updateTransaction
func (repo *transactionRepository) updateTransaction(t *Transaction) error {
	_, err := repo.db.Exec("UPDATE tblTransaction SET unitPrice = ?, totalAmount = ?, taxes = ? WHERE id = ?", t.unitPrice, t.totalAmount, t.taxes, t.id)
	return err
}

// deleteTransaction
func (repo *transactionRepository) deleteTransaction(id int) error {
	_, err := repo.db.Exec("DELETE FROM tblTransaction WHERE id = ?", id)
	return err
}
