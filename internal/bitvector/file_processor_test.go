package bitvector_test

import (
	"strings"
	"testing"

	"github.com/paulheg/kit_advanced_data_structures/internal/bitvector"
	"github.com/stretchr/testify/assert"
)

func TestFileProcessor(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc: "standard example",
			input: `6
001110110101010111111111
access 4
rank 0 10
select 1 14
rank 1 10
select 0 3
access 5`,
			expected: `1
4
20
6
5
0
`,
		},
		{
			desc: "simple test zeros",
			input: `7
000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
access 0
access 1
access 2
access 3
access 4
access 5
access 88
`,
			expected: `0
0
0
0
0
0
0
`,
		},
		{
			desc: "simple test access",
			input: `7
11110000111100001111000011110000111100001111000011110000111100001111000011110000111100001111000011110000111100001111000011110000
access 0
access 1
access 2
access 3
access 4
access 5
access 88
`,
			expected: `1
1
1
1
0
0
1
`,
		},
		{
			desc: "simple test access alternating",
			input: `5
0101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010
access 0
access 1
access 22
access 23
access 40
`,
			expected: `0
1
0
1
0
`,
		},
		{
			desc: "nothing working bro",
			input: `10
11101101111100111101100111101110110001111101011110001001101001011111000011001011000010010110100011000001101011110010011100000001011101100101111010010001100010111010110001000010000111111110100001001110111001111100010101001101101001110001010010000100000110100111100000001010001111110101111011011010101101110000000111010011000000001100000110111010100110011111110110000000111101110110001110010110010101011100001100001110001111111110111110011111111010010010101110010110011111100010111010011110011101101010110100100100110101011110010011110101111010011011010110111010100010001000001011111011010110011100111001010011011100100100110000011011011001000101001011010101101011111010011101100011111100100001100010100100101001001001001110010001000000100010011110000111001100011010110100101000000001101001011001100011001001101010000000000110101000111000010011010010111011100100100010101010100110011000100010100011000000010101001111101110000001110011001110001111100100000110011000111111100111110000001101110010011100101110010001011001011111100100100110000011101001110001011001011010101111100011110111001010011011011001001000000001011110110010011010000011100000011010000111000100110010011000010111100110110110110010010000010110101001000101111100000111010111011001111110010110111100100100100001001000
rank 0 698
access 1028
rank 0 284
access 1129
access 884
select 0 53
access 1228
select 1 65
access 13
select 1 73
`,
			expected: `329
1
139
0
1
120
0
114
0
133
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			input := strings.NewReader(tC.input)

			var output strings.Builder
			var statOut strings.Builder

			err := bitvector.ProcessFile(input, &output, &statOut, false)
			assert.NoError(t, err)

			assert.Equal(t, tC.expected, output.String())
			t.Log(statOut.String())
		})
	}
}
