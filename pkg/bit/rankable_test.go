package bit_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
	"github.com/stretchr/testify/assert"
)

type rankStragyBuilder func(bit.Vector) bit.Rankable

var rankStrategies = map[string]rankStragyBuilder{
	"baseline": func(v bit.Vector) bit.Rankable {
		return &bit.RankableBaseline{v}
	},
	"interleaved": func(v bit.Vector) bit.Rankable {
		return bit.NewInterleavedVector(v)
	},
}

func convert(input []byte) bit.Vector {
	vecLen := int(math.Ceil(float64(len(input)) / float64(bit.SubvectorBits)))

	vec := make(bit.Vector, vecLen)

	for i := 0; i < len(input); i++ {
		// pos := len(input) - i - 1
		if input[i] == 1 {
			vec.Set(uint64(i))
		}
	}

	return vec
}

func TestRankable(t *testing.T) {

	exampleVector := convert([]byte{0, 1, 1, 0, 1, 1, 0, 1, 0, 0})

	testCases := []struct {
		desc     string
		vec      bit.Vector
		pos      uint64
		alpha    bool
		expected uint64
	}{
		{
			desc:     "basic example zeros 5",
			vec:      exampleVector,
			pos:      5,
			alpha:    false,
			expected: 2,
		},
		{
			desc:     "basic example zeros 3",
			vec:      exampleVector,
			pos:      3,
			alpha:    false,
			expected: 1,
		},
		{
			desc:     "basic example ones",
			vec:      exampleVector,
			pos:      5,
			alpha:    true,
			expected: 3,
		},
		{
			desc:     "zero position does not include current position (zeros)",
			vec:      exampleVector,
			pos:      0,
			alpha:    false,
			expected: 0,
		},
		{
			desc:     "zero position does not include current position (ones)",
			vec:      exampleVector,
			pos:      0,
			alpha:    true,
			expected: 0,
		},
	}

	for stratName, strat := range rankStrategies {
		for _, tC := range testCases {
			t.Run(fmt.Sprintf("%s: %s", stratName, tC.desc), func(t *testing.T) {
				t.Logf("%b pos %d", tC.vec, tC.pos)
				r := strat(tC.vec)
				actual := r.Rank(tC.alpha, tC.pos)
				assert.Equal(t, tC.expected, actual)
			})
		}
	}
}

func BenchmarkRank(b *testing.B) {
	// generate vector
	const size = 8388608
	vector := make(bit.Vector, size)
	for i := 0; i < size; i++ {
		vector[i] = bit.Subvector(rand.Uint64())
	}
	b.ResetTimer()

	for stratName, strat := range rankStrategies {
		b.Run(fmt.Sprintf("%s", stratName), func(b *testing.B) {
			rv := strat(vector)

			for i := 0; i < 100; i++ {
				pos := rand.Int63n(int64(vector.Bits()))
				rv.Rank(true, uint64(pos))
			}
		})
	}
}

func TestRankVsNaive(t *testing.T) {
	// generate vector
	const size = 10_000
	vector := make(bit.Vector, size)
	for i := 0; i < size; i++ {
		vector[i] = bit.Subvector(rand.Uint64())
	}

	preparedVecs := make(map[string]bit.Rankable)
	for stratName, strat := range rankStrategies {
		preparedVecs[stratName] = strat(vector)
	}

	naive := bit.RankableBaseline{
		Vector: vector,
	}

	for i := 0; i < 1000; i++ {
		pos := rand.Int63n(size * int64(bit.SubvectorBits))
		expected1 := naive.Rank(true, uint64(pos))
		expected2 := naive.Rank(false, uint64(pos))

		for stratName, rankable := range preparedVecs {
			actual1 := rankable.Rank(true, uint64(pos))
			actual2 := rankable.Rank(false, uint64(pos))

			v1 := assert.Equal(t, expected1, actual1, stratName, expected1-actual1)
			v2 := assert.Equal(t, expected2, actual2, stratName, expected2-actual2)
			if !v1 || !v2 {
				continue
			}
		}
	}
}
