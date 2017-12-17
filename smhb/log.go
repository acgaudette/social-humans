package smhb

import (
	"bufio"
	"fmt"
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
		transaction.GetDir(),
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
		context.dataPath+"/transactions.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Fprintf(file, "%s\n", transaction.Timestamp)

	return nil
}

func countTransactions(context ServerContext) (int, error) {
	transactionLog.Lock()
	defer transactionLog.Unlock()

	file, err := os.Open(context.dataPath + "/transactions.log")

	if err != nil {
		log.Printf("countTransactions: %s", err.Error())
		return -1, err
	}

	defer file.Close()

	fs := bufio.NewScanner(file)
	lines := 0

	for fs.Scan() {
		lines++
	}

	log.Printf(
		"Log for %s:%d contains %d transactions",
		context.address, context.port, lines,
	)

	return lines, nil
}
