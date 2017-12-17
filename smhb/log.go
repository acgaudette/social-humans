package smhb

import (
	"bufio"
	"log"
	"os"
	"sync"
)

var transactionLog sync.Mutex

func (this *Transaction) GetDir() string {
	return "log/"
}

func (this *Transaction) GetPath() string {
	return this.GetDir() + this.Timestamp + ".trans"
}

func (this *Transaction) String() string {
	return "transaction \"" + this.Timestamp + "\""
}

/* Satisfy binary interfaces */

func (this *Transaction) MarshalBinary() ([]byte, error) {
	return writeTransaction(this)
}

func (this *Transaction) UnmarshalBinary(buffer []byte) error {
	this, err := readTransaction(buffer)
	return err
}

func logTransaction(
	transaction *Transaction, access Access, context ServerContext,
) error {
	err := access.SaveWithDir(
		transaction,
		context.dataPath+transaction.GetDir(),
		false,
		context,
	)

	if err != nil {
		return err
	}

	transactionLog.Lock()
	defer transactionLog.Unlock()

	// Write log to disk
	file, err := os.OpenFile(
		context.dataPath+transaction.GetDir()+"transactions.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func countTransactions(context ServerContext) (int, error) {
	transactionLog.Lock()
	defer transactionLog.Unlock()

	file, err := os.Open(context.dataPath + "/transactions/transactions.log")
	defer file.Close()
	if err != nil {
		log.Printf("countTransactions: %s", err.Error())
		return 0, err
	}

	fs := bufio.NewScanner(file)
	lines := 0

	for fs.Scan() {
		lines++
	}

	log.Printf(
		"countTransactions: log for %s:%d contains %d transactions",
		context.address, context.port, lines,
	)

	return lines, nil
}
