package bit_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
	"github.com/stretchr/testify/assert"
)

type selectableStrategyBulder func(vec bit.Vector) bit.Selectable

var selectStrategies = map[string]selectableStrategyBulder{
	"baseline": func(vec bit.Vector) bit.Selectable {
		return &bit.SelectableBaseline{
			Vector: vec,
		}
	},
	"interleaved": func(vec bit.Vector) bit.Selectable {
		return bit.NewInterleavedVector(vec)
	},
}

func BenchmarkSelect(b *testing.B) {
	const vecSize = 10000000
	const operations = 100

	// generate random vector
	vector := make(bit.Vector, vecSize)
	for i := 0; i < vecSize; i++ {
		vector[i] = bit.Subvector(rand.Uint64())
	}
	ones := vector.Ones()

	b.ResetTimer()

	for stratName, strat := range selectStrategies {
		b.Run(fmt.Sprintf("%s", stratName), func(b *testing.B) {

			// build time
			rv := strat(vector)

			// selects
			for i := 0; i < operations; i++ {
				pos := uint64(rand.Int63n(int64(ones)) + 1)
				rv.Select(true, pos)
			}
		})
	}
}

func TestSelectVsNaive(t *testing.T) {
	// generate vector
	const size = 10_000
	vector := make(bit.Vector, size)
	for i := 0; i < size; i++ {
		vector[i] = bit.Subvector(rand.Uint64())
	}

	ones := vector.Ones()
	zeros := vector.Bits() - ones

	preparedVecs := make(map[string]bit.Selectable)
	for stratName, strat := range selectStrategies {
		preparedVecs[stratName] = strat(vector)
	}

	naive := bit.SelectableBaseline{
		Vector: vector,
	}

	for i := 0; i < 1000; i++ {
		onePos := rand.Int63n(int64(ones))
		zeroPos := rand.Int63n(int64(zeros))

		expected1 := naive.Select(true, uint64(onePos))
		expected2 := naive.Select(false, uint64(zeroPos))

		for stratName, rankable := range preparedVecs {
			actual1 := rankable.Select(true, uint64(onePos))
			actual2 := rankable.Select(false, uint64(zeroPos))

			assert.Equal(t, expected1, actual1, stratName)
			assert.Equal(t, expected2, actual2, stratName)
		}
	}
}

func TestSelectablePanics(t *testing.T) {
	testCases := []struct {
		desc  string
		vec   bit.Vector
		alpha bool
		n     uint64
	}{
		{
			desc:  "panic when using 0 bit position",
			vec:   convert([]byte{0, 1, 0, 1, 0}),
			alpha: true,
			n:     0,
		},
		{
			desc:  "panic when using 0 bit position",
			vec:   convert([]byte{0, 1, 0, 1, 0}),
			alpha: false,
			n:     0,
		},
	}
	for stratName, strat := range selectStrategies {
		for _, tC := range testCases {
			t.Run(fmt.Sprintf("%s: %s", stratName, tC.desc), func(t *testing.T) {
				s := strat(tC.vec)
				assert.Panics(t, func() {
					s.Select(tC.alpha, tC.n)
				})
			})
		}
	}
}

func TestSelectable(t *testing.T) {

	exampleVector := convert([]byte{0, 1, 1, 0, 1, 1, 0, 1, 0, 0})
	//								0  1  2  3  4  5  6  7  8  9

	bigVector := convert([]byte{
		0, 1, 0, 1, 0, 0, 0, 1, 1, 0, // 0   6#0
		1, 0, 1, 0, 0, 0, 0, 0, 1, 1, // 10  12#0
		1, 1, 0, 1, 0, 0, 0, 0, 0, 0, // 20  19#0
		1, 1, 0, 1, 0, 1, 0, 0, 1, 1, // 30  23#0
		1, 1, 1, 0, 0, 0, 1, 1, 1, 0, // 40  27#0
		1, 1, 0, 1, 0, 1, 0, 1, 0, 0, // 50  32#0
		1, 0, 0, 0, 0, 1, 1, 1, 0, 0, // 60
	})
	//  0  1  2  3  4  5  6  7  8  9

	otherBigVector := bit.Vector{
		bit.Subvector(0xFFFFFFFF_EEEEEEEE), // 8 #0
		bit.Subvector(0xFFFFFFFE_FFFFFFFF), // 9 #0
	}

	testCases := []struct {
		desc     string
		vec      bit.Vector
		n        uint64
		alpha    bool
		expected uint64
	}{
		{
			desc:     "basic example zeros",
			vec:      exampleVector,
			n:        5,
			alpha:    false,
			expected: 9,
		},
		{
			desc:     "basic example zeros 2",
			vec:      exampleVector,
			n:        2,
			alpha:    false,
			expected: 3,
		},
		{
			desc:     "basic example zeros",
			vec:      exampleVector,
			n:        3,
			alpha:    false,
			expected: 6,
		},
		{
			desc:     "basic example ones",
			vec:      exampleVector,
			n:        5,
			alpha:    true,
			expected: 7,
		},
		{
			desc:     "big vector zeros over 64",
			vec:      bigVector,
			n:        34,
			expected: 62,
			alpha:    false,
		},
		{
			desc:     "other big vector zeros under 64",
			vec:      otherBigVector,
			n:        2,
			expected: 4,
			alpha:    false,
		},
		{
			desc:     "other big vector zeros over 64",
			vec:      otherBigVector,
			n:        9,
			expected: 96,
			alpha:    false,
		},
		// {
		// 	desc:     "big vector zeros over 64",
		// 	vec:      bigVector,
		// 	n:        36,
		// 	expected: 64,
		// 	alpha:    false,
		// },
		// {
		// 	desc:     "big vector zeros over 64",
		// 	vec:      bigVector,
		// 	n:        37,
		// 	expected: 68,
		// 	alpha:    false,
		// },
	}

	for stratName, strat := range selectStrategies {
		for _, tC := range testCases {
			t.Run(fmt.Sprintf("%s: %s", stratName, tC.desc), func(t *testing.T) {
				t.Logf("%b pos %d alpha %v", tC.vec, tC.n, tC.alpha)
				s := strat(tC.vec)
				actual := s.Select(tC.alpha, tC.n)
				assert.Equal(t, tC.expected, actual)
			})
		}
	}
}
