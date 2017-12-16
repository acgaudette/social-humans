package smhb

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var transactionLog sync.Mutex

func (this *Transaction) GetDir() string {
	return "transactions/"
}

func (this *Transaction) GetPath() string {
	return this.GetDir() + this.Timestamp
}

func (this *Transaction) String() string {
	return this.Timestamp
}

func (this *Transaction) MarshalBinary() ([]byte, error) {
	return serialize(*this)
}

func (this *Transaction) UnmarshalBinary(buffer []byte) error {
	err := deserialize(this, buffer)
	if err != nil {
		return err
	}
	return nil
}

func logTransaction(
	transaction *Transaction, access Access, context ServerContext,
) error {
	access.SaveWithDir(transaction, transaction.GetDir(), false, context)

	transactionLog.Lock()
	defer transactionLog.Unlock()

	file, err := os.OpenFile("transactions/transactions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	fmt.Fprintf(file, "%s\n", transaction.Timestamp)

	return nil
}

func countTransactions() (int, error) {
	transactionLog.Lock()
	defer transactionLog.Unlock()

	file, err := os.Open("transactions/transaction.log")
	defer file.Close()
	if err != nil {
		return 0, err
	}

	fs := bufio.NewScanner(file)
	lines := 0
	for fs.Scan() {
		lines++
	}
	return lines, nil
}
