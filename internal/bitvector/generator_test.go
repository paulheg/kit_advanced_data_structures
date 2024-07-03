package bitvector_test

import (
	"strings"
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/internal/bitvector"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {

	var commands strings.Builder
	var expected strings.Builder

	err := bitvector.GenerateRandomTestCase(10, 1000, &commands, &expected)
	assert.NoError(t, err)
}
