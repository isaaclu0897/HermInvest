package service

import (
	"HermInvest/pkg/model"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

type service struct {
	repo model.Repositorier
}

func NewService(repository model.Repositorier) *service {
	return &service{repo: repository}
}

func (serv *service) WithTrx(trxHandle *gorm.DB) *service {
	return &service{repo: serv.repo.WithTrx(trxHandle)} // return new one
}

// addTransactionTailRecursion add new transaction records with tail recursion,
// When adding, inventory and transaction history, especially write-offs and
// tails, need to be considered.
func (serv *service) addTransactionTailRecursion(newTransaction *model.Transaction, remainingQuantity int) (*model.Transaction, error) {
	// Principles:
	// 1. Ensure that each transaction has a corresponding transaction record.
	// 2. Update inventory quantities based on transactions, including adding,
	//    reducing, or deleting inventory.
	// 3. Depending on the transaction situation, only transaction history can
	//    be added and cannot be modified or deleted.
	// 4. For insufficient write-off quantities, recursive processing is used
	//    to ensure that the write-off is completed.

	// Cases:
	// 1. Newly added: If there is no transaction in the inventory (A) or
	//    the new transaction is the same as the oldest transaction in the
	//    inventory (B), add it directly to the inventory.
	// 2. Write-off:
	// 	* Sufficient inventory: If the inventory quantity is sufficient,
	//    update the inventory quantity (C) or delete the inventory (D), and
	//    add the corresponding transaction history.
	// 	* Insufficient inventory: If the inventory quantity can't be Write-off.
	//    Recurse until success (E). The termination condition is A B C D.
	//  * Over inventory: Write-off over than inventory (F).

	// TODO: This func should be moved to service tier.

	earliestTransaction, err := serv.repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to finding first purchase: %v", err)
		}
		// Case A
		earliestTransaction.TranType = newTransaction.TranType
	}

	if earliestTransaction.TranType == newTransaction.TranType {
		if newTransaction.Quantity != remainingQuantity {
			// Case F
			newTransaction.SetQuantity(newTransaction.Quantity - remainingQuantity)
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(F), failed to creating transaction history: %v", err)
			}
			newTransaction.SetQuantity(remainingQuantity)
		}

		// Case B
		id, err := serv.repo.CreateTransaction(newTransaction)
		if err != nil {
			return nil, fmt.Errorf("Case(B), failed to creating transaction: %v", err)
		}
		transaction, err := serv.repo.QueryTransactionByID(id)
		if err != nil {
			return nil, fmt.Errorf("Case(B), failed to querying transaction: %v", err)
		}

		return transaction, nil
	} else {
		if earliestTransaction.Quantity > remainingQuantity {
			// Case C

			// Create a copy for adding stock history
			stockHistoryAdd := &model.Transaction{}
			*stockHistoryAdd = *earliestTransaction
			// var stockHistoryAdd *model.Transaction // why can't use it, study it
			// *stockHistoryAdd = *earliestTransaction

			// add transaction history
			stockHistoryAdd.SetQuantity(remainingQuantity)
			_, err = serv.repo.CreateTransactionHistory(stockHistoryAdd)
			if err != nil {
				return nil, fmt.Errorf("Case(C), failed to creating transaction history: %v", err)
			}
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(C), failed to creating transaction history: %v", err)
			}

			// Update stock inventory
			earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
			err := serv.repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(C), failed to updating transaction: %v", err)
			}

			return earliestTransaction, nil
		} else if earliestTransaction.Quantity == remainingQuantity {
			// Case D

			// add transaction history
			_, err = serv.repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(D), failed to creating transaction history: %v", err)
			}
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(D), failed to creating transaction history: %v", err)
			}
			// delete stock inventory
			err = serv.repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				return nil, fmt.Errorf("Case(D), failed to deleting transaction: %v", err)
			}

			// Or use move

			return nil, nil
		} else { // earliestTransaction.Quantity < remainingQuantity
			// Case E

			// add transaction history
			_, err = serv.repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				return nil, fmt.Errorf("Case(E), failed to creating transaction history: %v", err)
			}

			// delete stock inventory
			err = serv.repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				return nil, fmt.Errorf("Case(E), failed to deleting transaction: %v", err)
			}

			remainingQuantity = remainingQuantity - earliestTransaction.Quantity

			return serv.addTransactionTailRecursion(newTransaction, remainingQuantity)
		}
	}
}

// AddTransaction add the transaction from the input to the inventory.
// It will add or update transactions in the inventory and add history.
// Return the modified transaction record in the inventory
func (serv *service) AddTransaction(newTransaction *model.Transaction) (*model.Transaction, error) {
	tx := serv.repo.Begin()

	remainingQuantity := newTransaction.Quantity
	ts, err := serv.WithTrx(tx).addTransactionTailRecursion(newTransaction, remainingQuantity)
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return nil, fmt.Errorf("failed to add transaction: %v", err)
	}
	serv.repo.WithTrx(tx).Commit()

	return ts, nil
}

// ---

func (serv *service) DeleteTransaction(id int) error {
	return serv.repo.DeleteTransaction(id)
}

func (serv *service) QueryTransactionAll() ([]*model.Transaction, error) {
	return serv.repo.QueryTransactionAll()
}

func (serv *service) QueryTransactionByID(id int) (*model.Transaction, error) {
	return serv.repo.QueryTransactionByID(id)
}

func (serv *service) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
	return serv.repo.QueryTransactionByDetails(stockNo, tranType, date)
}

func (serv *service) UpdateTransaction(id int, t *model.Transaction) error {
	return serv.repo.UpdateTransaction(id, t)
}

// ---

type DividendOrReduction struct {
	Date string
	Obj  interface{}
}

func mergeAndSort(exDividends []*model.ExDividend, capitalReductions []*model.CapitalReduction) []*DividendOrReduction {
	var mergedList []*DividendOrReduction

	for _, exDividend := range exDividends {
		mergedList = append(mergedList, &DividendOrReduction{Date: exDividend.ExDividendDate, Obj: exDividend})
	}

	for _, capitalReduction := range capitalReductions {
		mergedList = append(mergedList, &DividendOrReduction{Date: capitalReduction.CapitalReductionDate, Obj: capitalReduction})
	}

	sort.Slice(mergedList, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02", mergedList[i].Date)
		date2, _ := time.Parse("2006-01-02", mergedList[j].Date)
		return date1.Before(date2)
	})

	return mergedList
}

func (serv *service) RebuildTransactionRecordSys() error {

	eds, err := serv.repo.QueryDividendAll()
	if err != nil {
		return err
	}

	crs, err := serv.repo.QueryCapitalReductionAll()
	if err != nil {
		return err
	}

	trs, err := serv.repo.QueryTransactionRecordAll()
	if err != nil {
		return err
	}

	mergedList := mergeAndSort(eds, crs)

	var cashDividends []*model.ExDividend
	for _, o := range mergedList {
		var filteredRecords []*model.TransactionRecord

		switch obj := o.Obj.(type) {
		case *model.CapitalReduction:
			cr := obj
			for _, record := range trs {
				rdate, _ := time.Parse("2006-01-02", record.Date)
				crdate, _ := time.Parse("2006-01-02", cr.CapitalReductionDate)
				if cr.StockNo == record.StockNo && crdate.After(rdate) {
					filteredRecords = append(filteredRecords, record)
				}
			}

			remainingTrs, err := model.CalcRemainingTransactionRecords(filteredRecords)
			if err != nil {
				return err
			}

			totalQuantity, avgUnitPrice := model.SumQuantityUnitPrice(remainingTrs)

			capitalReductionRecord, distributionRecord := cr.CalcTransactionRecords(totalQuantity, avgUnitPrice)

			trs = append(trs, capitalReductionRecord, distributionRecord)
		case *model.ExDividend:
			ed := obj
			for _, record := range trs {
				rdate, _ := time.Parse("2006-01-02", record.Date)
				crdate, _ := time.Parse("2006-01-02", ed.ExDividendDate)
				if ed.StockNo == record.StockNo && crdate.After(rdate) {
					filteredRecords = append(filteredRecords, record)
				}
			}

			remainingTrs, err := model.CalcRemainingTransactionRecords(filteredRecords)
			if err != nil {
				return err
			}

			totalQuantity, _ := model.SumQuantityUnitPrice(remainingTrs)

			// TODO: need to calc stock dividend record and append to newTrs
			cd := ed.CalcCashDividendRecord(totalQuantity)

			cashDividends = append(cashDividends, cd)
		}

		sort.Slice(trs, func(i, j int) bool {
			date1, _ := time.Parse("2006-01-02", trs[i].Date)
			date2, _ := time.Parse("2006-01-02", trs[j].Date)
			return date1.Before(date2)
		})
	}

	tx := serv.repo.Begin()

	err = serv.repo.WithTrx(tx).DropTable("tblTransactionCash")
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return err
	}

	for _, cd := range cashDividends {
		err = serv.repo.WithTrx(tx).CreateCashDividendRecord(cd)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return err
		}
	}

	err = serv.repo.WithTrx(tx).DropTable("tblTransactionRecordSys")
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return err
	}

	for _, tr := range trs {
		err = serv.repo.WithTrx(tx).CreateTransactionRecordSys(tr)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return err
		}
	}

	serv.repo.WithTrx(tx).Commit()

	return nil
}

func (serv *service) RebuildTransaction() error {
	tx := serv.repo.Begin()

	err := serv.repo.WithTrx(tx).DropTable("sqlite_sequence")
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return fmt.Errorf("failed to deleting SQLiteSequence: %v", err)
	}

	err = serv.repo.WithTrx(tx).DropTable("tblTransaction")
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return fmt.Errorf("failed to deleting tblTransaction: %v", err)
	}

	err = serv.repo.WithTrx(tx).DropTable("tblTransactionHistory")
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return fmt.Errorf("failed to deleting tblTransactionHistory: %v", err)
	}

	trs, err := serv.repo.WithTrx(tx).QueryTransactionRecordSysAll()
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return fmt.Errorf("failed to querying TransactionRecord: %v", err)
	}

	for _, tr := range trs {
		newTransaction := model.NewTransactionFromInput(
			tr.Date, tr.Time, tr.StockNo, tr.TranType, tr.Quantity, tr.UnitPrice)
		remainingQuantity := newTransaction.Quantity
		_, err := serv.WithTrx(tx).addTransactionTailRecursion(newTransaction, remainingQuantity)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return fmt.Errorf("failed to adding transaction in tail recursion: %v", err)
		}
	}

	serv.repo.WithTrx(tx).Commit()

	return nil
}
