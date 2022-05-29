package csv

import (
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/transactions"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	t.Run("Validate output conversion to CSV", func(t *testing.T) {
		mockTransaction := transactions.Build(
			time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)).GenerateTransaction()
		w := Writer{
			data: []transactions.Transaction{
				mockTransaction,
			},
		}
		expectedRecord := mockTransaction.ToRecord()
		csvBuffer, err := w.ToCsv()
		t.Log(string(csvBuffer))
		assert.ErrorIs(t, err, nil)
		assert.Contains(t, string(csvBuffer), strings.Join(expectedRecord, ","))
	})
}
