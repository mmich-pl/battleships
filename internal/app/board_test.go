package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapCoords(t *testing.T) {
	testScenarios := []struct {
		name   string
		coords string
		error  bool
	}{
		{"valid coords - one int coords", "A1", false},
		{"valid coords - two int coords", "B10", false},
		{"invalid coords - two letter coords", "AB1", true},
		{"invalid coords - negative int", "A-1", true},
		{"invalid coords - letter out of bounds ", "Z1", true},
		{"invalid coords - int out of bounds", "A13", true},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, _, err := mapCoords(scenario.coords)
			if scenario.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
