package transactions

import (
	"fmt"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/dates"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const MaxValue = 100000

type RandomDateGenerator interface {
	GenerateRandomDate() dates.RandomDate
}

type TransactionGenerator struct {
	dateGenerator RandomDateGenerator
}

type TransactionGeneratorInterface interface {
	GenerateTransaction() Transaction
	GenerateTransactionWithError() Transaction
}

type Transaction struct {
	Id     string
	Date   dates.RandomDate
	Amount float64
}

func Build(initialDate time.Time, finalDate time.Time) *TransactionGenerator {
	randomDataGenerator := dates.Build(initialDate, finalDate)
	t := TransactionGenerator{dateGenerator: randomDataGenerator}
	return &t
}

func (t TransactionGenerator) GenerateTransaction() Transaction {
	rand.Seed(time.Now().UnixNano())
	return Transaction{
		Id:     uuid.New().String(),
		Date:   t.dateGenerator.GenerateRandomDate(),
		Amount: (rand.Float64() - 0.5) * MaxValue,
	}
}

func (t TransactionGenerator) GenerateTransactionWithError() Transaction {
	rand.Seed(time.Now().UnixNano())
	return Transaction{
		Id:     "invalid",
		Date:   t.dateGenerator.GenerateRandomDate(),
		Amount: 0.0,
	}
}

func (t Transaction) ToRecord() []string {
	record := make([]string, 3)
	record[0] = t.Id
	record[1] = t.Date.Date.String()
	record[2] = fmt.Sprintf("%0.2f", t.Amount)
	return record
}
