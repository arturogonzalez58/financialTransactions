package transactions

import (
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/dates"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type DateGeneratorMock struct {
	isCalled     bool
	mockResponse dates.RandomDate
}

func (d DateGeneratorMock) GenerateRandomDate() dates.RandomDate {
	d.isCalled = true
	return d.mockResponse
}

func TestTransactionGenerator(t1 *testing.T) {
	t1.Run("Validate transaction", func(t1 *testing.T) {
		initialDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
		finalDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
		dateGenerateMock := DateGeneratorMock{
			mockResponse: dates.Build(initialDate, finalDate).GenerateRandomDate(),
		}
		t := TransactionGenerator{
			dateGenerator: dateGenerateMock,
		}
		transaction := t.GenerateTransaction()
		assert.Equalf(t1, transaction.Date, dateGenerateMock.mockResponse, "Should be equal")
		assert.NotEqual(t1, transaction.Amount, 0)
		assert.NotEqual(t1, transaction.Id, "")
	})

	t1.Run("Validate transaction with error", func(t1 *testing.T) {
		initialDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
		finalDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
		dateGenerateMock := DateGeneratorMock{
			mockResponse: dates.Build(initialDate, finalDate).GenerateRandomDate(),
		}
		t := TransactionGenerator{
			dateGenerator: dateGenerateMock,
		}
		transaction := t.GenerateTransactionWithError()
		assert.Equal(t1, transaction.Date, dateGenerateMock.mockResponse, "Should be equal")
		assert.Equal(t1, transaction.Amount, 0.00)
		assert.Equal(t1, transaction.Id, "")
	})

	t1.Run("Transaction to record", func(t1 *testing.T) {
		initialDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
		finalDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
		dateGenerateMock := DateGeneratorMock{
			mockResponse: dates.Build(initialDate, finalDate).GenerateRandomDate(),
		}
		t := TransactionGenerator{
			dateGenerator: dateGenerateMock,
		}
		transaction := t.GenerateTransactionWithError()
		record := transaction.ToRecord()
		assert.Equal(t1, record[1], dateGenerateMock.mockResponse.Date.String(), "Should be equal")
		assert.Equal(t1, record[2], "0.000000")
		assert.Equal(t1, record[0], "")
	})
}
