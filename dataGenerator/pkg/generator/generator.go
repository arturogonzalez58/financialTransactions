package generator

import (
	"fmt"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/transactions"
	"math/rand"
	"time"
)

type TransactionGenerator interface {
	GenerateTransaction() transactions.Transaction
	GenerateTransactionWithError() transactions.Transaction
}

type Generator struct {
	transactionGenerator TransactionGenerator
	dataSize             int32
	errorPercentage      float32
}

func Builder(dataSie int32, errorPercentage float32, initialDate time.Time, finalDate time.Time) *Generator {
	transactionsGenerator := transactions.Build(initialDate, finalDate)
	g := Generator{
		dataSize:             dataSie,
		errorPercentage:      errorPercentage,
		transactionGenerator: transactionsGenerator,
	}
	return &g
}

func (g Generator) GenerateData() []transactions.Transaction {
	data := make([]transactions.Transaction, g.dataSize)
	for i := range data {
		if !g.isAnError() {
			data[i] = g.transactionGenerator.GenerateTransaction()
		} else {
			data[i] = g.transactionGenerator.GenerateTransactionWithError()
		}
	}
	return data
}

func (g Generator) isAnError() bool {
	fmt.Printf("%f", g.errorPercentage)
	fmt.Printf("%f", rand.Float32())
	return g.errorPercentage > rand.Float32()
}
