package generator

import (
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/dates"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/transactions"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TransactionGeneratorMock struct {
	callGenerate          bool
	responseGenerate      transactions.Transaction
	callGenerateError     bool
	responseGenerateError transactions.Transaction
}

func (t TransactionGeneratorMock) GenerateTransaction() transactions.Transaction {
	t.callGenerate = true
	return t.responseGenerate
}

func (t TransactionGeneratorMock) GenerateTransactionWithError() transactions.Transaction {
	t.callGenerateError = true
	return t.responseGenerateError
}

func TestGenerator(t *testing.T) {
	t.Run("Generate data", func(t *testing.T) {

		mockDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
		mockTransaction := transactions.Transaction{Id: "valid-id", Date: dates.RandomDate{Date: mockDate}, Amount: 55.55}
		mockTransactionError := transactions.Transaction{Id: "invalid-id", Date: dates.RandomDate{Date: mockDate}, Amount: 0}
		transactionGeneratorMock := TransactionGeneratorMock{
			responseGenerate:      mockTransaction,
			responseGenerateError: mockTransactionError,
		}
		g := Generator{
			transactionGenerator: transactionGeneratorMock,
			dataSize:             10,
			errorPercentage:      0.10,
		}
		data := g.GenerateData()
		assert.True(t, len(data) == 10)
		assert.Contains(t, data, mockTransaction)
	})
}
