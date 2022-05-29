package dates

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRandomDateGenerator_GenerateRandomDate(t *testing.T) {

	t.Run("Generate a random date between the initial date and the final date", func(t *testing.T) {
		initialDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
		finalDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
		g := Build(initialDate, finalDate)
		randomDate := g.GenerateRandomDate()
		assert.True(t, randomDate.Date.After(initialDate) && randomDate.Date.Before(finalDate), ""+
			"The random date must be less than final date and more than initial date")
	})
}
