package dates

import (
	"math/rand"
	"time"
)

type RandomDateGenerator struct {
	initialDate time.Time
	finalDate   time.Time
}

type RandomDate struct {
	Date time.Time
}

func Build(initialDate time.Time, finalDate time.Time) *RandomDateGenerator {
	r := RandomDateGenerator{initialDate: initialDate, finalDate: finalDate}
	return &r
}

func (g RandomDateGenerator) GenerateRandomDate() RandomDate {
	delta := g.finalDate.Unix() - g.initialDate.Unix()
	sec := rand.Int63n(delta) + g.initialDate.Unix()
	return RandomDate{Date: time.Unix(sec, 0)}
}
