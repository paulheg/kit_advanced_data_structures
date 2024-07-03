package bit_test

import (
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
	"github.com/stretchr/testify/assert"
)

func TestInterleavedRank(t *testing.T) {

	onesVec := bit.Vector{
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 64
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 128
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 192
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 256
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 320
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 384
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 448
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 512
		bit.Subvector(0xFFFFFFFF_FFFFFFFF), // 576
	}

	testCases := []struct {
		desc     string
		vec      bit.Vector
		alpha    bool
		position uint64
		expected uint64
	}{
		{
			desc:     "zeros in one vec - pos 266",
			vec:      onesVec,
			alpha:    false,
			position: 266,
			expected: 0,
		},
		{
			desc:     "zeros in one vec - pos 522 - second line",
			vec:      onesVec,
			alpha:    false,
			position: 522,
			expected: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interleaved := bit.NewInterleavedVector(tC.vec)
			actual := interleaved.Rank(tC.alpha, tC.position)

			assert.Equal(t, tC.expected, actual)
		})
	}
}
