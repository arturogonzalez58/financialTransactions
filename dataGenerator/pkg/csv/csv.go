package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/transactions"
)

type Writer struct {
	data []transactions.Transaction
}

func Build(data []transactions.Transaction) *Writer {
	writer := Writer{data: data}
	return &writer
}

func (w Writer) ToCsv() ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for _, transaction := range w.data {
		record := transaction.ToRecord()
		err := writer.Write(record)
		writer.Flush()
		if err != nil {
			return nil, fmt.Errorf("there was an error creating the csv: %w", err)
		}
	}
	return buf.Bytes(), nil
}
