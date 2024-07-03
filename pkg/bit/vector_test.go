package bit_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
	"github.com/stretchr/testify/assert"
)

func TestBitVectorSetAndGet(t *testing.T) {
	testCases := []struct {
		desc     string
		position uint64
		size     uint64
	}{
		{
			desc:     "position 0",
			position: 0,
			size:     1,
		},
		{
			desc:     "position 1",
			position: 1,
			size:     1,
		},
		{
			desc:     "position 63",
			position: 63,
			size:     1,
		},
		{
			desc:     "bigger: position 64",
			position: 64,
			size:     2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			bv := make(bit.Vector, tC.size)

			bv.Set(tC.position)
			t.Logf("%b", bv)
			assert.True(t, bv.Access(tC.position))

			bv.Unset(tC.position)
			assert.False(t, bv.Access(tC.position))
		})
	}
}

func TestSubvector(t *testing.T) {
	vect := bit.Vector{0xFF00FF00FF00FF00, 0xFF00FF00FF00FF00}

	testCases := []struct {
		desc             string
		vec              bit.Vector
		position, length uint64
		expected         bit.Subvector
	}{
		{
			desc:     "begin",
			position: 0,
			length:   8,
			expected: 0x00,
			vec:      vect,
		},
		{
			desc:     "begin with offset",
			position: 8,
			length:   8,
			expected: 0xFF,
			vec:      vect,
		},
		{
			desc:     "overlap",
			position: 60,
			length:   8,
			expected: 0x0F,
			vec:      vect,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			actual := tC.vec.Subvector(tC.position, tC.length)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestSubvectorSelect(t *testing.T) {
	testCases := []struct {
		desc     string
		vector   bit.Subvector
		pos      uint8
		alpha    bool
		expected uint8
	}{
		{
			desc:     "3. one",
			vector:   bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0111_0111),
			alpha:    true,
			pos:      3,
			expected: 2,
		},
		{
			desc:     "1. one",
			vector:   bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0000_0000),
			alpha:    true,
			pos:      1,
			expected: 8,
		},
		{
			desc:     "3. zero",
			vector:   bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0111_0111),
			alpha:    false,
			pos:      3,
			expected: 11,
		},
		{
			desc:     "starting 1",
			vector:   bit.Subvector(0x00_00_00_00),
			alpha:    false,
			pos:      5,
			expected: 4,
		},
		{
			desc:     "starting 1",
			vector:   bit.Subvector(0x00_00_00_00),
			alpha:    false,
			pos:      5,
			expected: 4,
		},
		{
			desc:     "",
			vector:   bit.Subvector(0b1011_0110),
			alpha:    false,
			pos:      4,
			expected: 8,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.vector.Select(tC.alpha, tC.pos)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestSubvectorSelectOne64(t *testing.T) {
	testCases := []struct {
		desc     string
		vector   bit.Subvector
		pos      uint8
		expected uint8
	}{
		{
			desc:     "starting 1",
			vector:   bit.Subvector(0x00_00_00_00_00_00_00_00_00_01),
			pos:      1,
			expected: 0,
		},
		{
			desc:     "only 1 10th",
			vector:   bit.Subvector(0xFF_FF_FF_FF),
			pos:      10,
			expected: 9,
		},
		{
			desc:     "only 1 9th",
			vector:   bit.Subvector(0xFF_FF_FF_FF),
			pos:      9,
			expected: 8,
		},
		{
			desc:     "only 1 7th",
			vector:   bit.Subvector(0xFF_FF_FF_FF),
			pos:      7,
			expected: 6,
		},
		{
			desc:     "",
			vector:   bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0111_0111),
			pos:      3,
			expected: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.vector.OneSelect64(tC.pos)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestOneSelect64_WithOnlyOnes(t *testing.T) {
	vec := bit.Subvector(0xFF_FF_FF_FF)

	for i := uint8(1); i <= 16; i++ {
		t.Run(fmt.Sprintf("Test pos %d", i), func(t *testing.T) {
			pos := vec.OneSelect64(i)
			assert.Equal(t, i-1, pos)
		})
	}
}

func BenchmarkSubvectorSelect(b *testing.B) {

	vec := make(bit.Vector, b.N)
	for i := range len(vec) {
		vec[i] = bit.SubvectorMax
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vec[i].Select(true, uint8(rand.Intn(64)+1))
	}
}

func TestSubvectorRank(t *testing.T) {
	testCases := []struct {
		desc     string
		input    bit.Subvector
		expected uint8
		alpha    bool
		n        uint8
	}{
		{
			desc:     "",
			input:    bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0000_0000),
			alpha:    true,
			expected: 0,
			n:        8,
		},
		{
			desc:     "",
			input:    bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0000_0000),
			alpha:    true,
			expected: 1,
			n:        9,
		},
		{
			desc:     "",
			input:    bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0000_0000),
			alpha:    false,
			expected: 8,
			n:        9,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.input.Rank(tC.alpha, tC.n))
		})
	}
}

func TestSubvectorAccess(t *testing.T) {
	subvec := bit.Subvector(0b_01111_0111_0111_0111_01111_0111_0000_0001)

	testCases := []struct {
		desc     string
		input    bit.Subvector
		expected bool
		pos      uint8
	}{
		{
			desc:     "pos 7",
			input:    subvec,
			expected: false,
			pos:      7,
		},
		{
			desc:     "pos 0",
			input:    subvec,
			expected: true,
			pos:      0,
		},
		{
			desc:     "pos 1",
			input:    subvec,
			expected: false,
			pos:      1,
		},
		{
			desc:     "pos 11",
			input:    subvec,
			expected: false,
			pos:      11,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.input.Access(tC.pos))
		})
	}
}
