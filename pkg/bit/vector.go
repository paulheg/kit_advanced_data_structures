package bit

import (
	"fmt"
	"math"
	"math/bits"
)

const SubvectorBits uint64 = 64
const SubvectorMax Subvector = math.MaxUint64

type Subvector uint64

func (s *Subvector) Set(pos uint8) {
	*s |= 1 << pos
}

func (s Subvector) Ones() uint8 {
	return uint8(bits.OnesCount64(uint64(s)))
}

func (s Subvector) Access(pos uint8) bool {
	return (s&(1<<pos))>>pos == 1
}

func (s Subvector) Rank(alpha bool, n uint8) uint8 {
	offsetBlock := s & ^(math.MaxUint64 << n)
	ones := uint8(bits.OnesCount64(uint64(offsetBlock)))
	if alpha {
		return ones
	} else {
		return n - ones
	}
}

func (s Subvector) Select(alpha bool, n uint8) uint8 {

	if alpha {
		return s.OneSelect64(n)
	} else {
		return (^s).OneSelect64(n)
	}
}

func (s Subvector) OneSelect64(n uint8) uint8 {

	const posMask = 0b0111
	const validMask = 0b1000

	n -= 1

	for i := 0; i < int(SubvectorBits); i += 8 {
		segment := uint8(s >> i)
		ones := uint8(bits.OnesCount8(segment))
		if ones < n+1 {
			n -= ones
		} else {

			// 4 bytes per number
			var lookupBegin uint32 = uint32(segment) * 4
			bytePos := lookupBegin + (uint32(n) >> 1)
			valueTuple := byte(onesLookup[bytePos])
			value := valueTuple >> ((n % 2) * 4)

			isValid := value&validMask == validMask
			if isValid {
				// add segment offset
				return (value & posMask) + byte(i)
			} else {
				panic("not found")
			}
		}
	}

	// should not happen
	panic("not found")
}

var _ Accessible = (*Vector)(nil)
var _ AccessibleWithSize = (*Vector)(nil)

type Vector []Subvector

// Size implements AccessibleWithSize.
func (b *Vector) Size() uint64 {
	return b.Bits()
}

type Accessible interface {
	Access(position uint64) bool
}

type Setable interface {
	Set(position uint64)
}

type AccessibleWithSize interface {
	Accessible
	Sizable
}

type RankSelectVector interface {
	AccessibleWithSize
	RankableWithSize
	SelectableWithSize

	// Additional bits
	Overhead() uint64
}

func NewVector(input string) Vector {
	var noSubvectors = uint64(math.Ceil(float64(len(input)) / float64(SubvectorBits)))
	var vec Vector = make([]Subvector, noSubvectors)

	for i := 0; i < len(input); i++ {
		// pos := len(input) - i - 1
		if input[i] == '1' {
			vec.Set(uint64(i))
		}
	}

	return vec
}

func (b Vector) Set(position uint64) {
	subvectorPos := position / SubvectorBits
	bitPos := position % SubvectorBits

	b[subvectorPos] |= 1 << bitPos
}

func (b Vector) Unset(position uint64) {
	subvectorPos := position / SubvectorBits
	bitPos := position % SubvectorBits

	b[subvectorPos] &= ^(1 << bitPos)
}

func (b Vector) Access(position uint64) bool {
	subvectorPos := position / SubvectorBits
	bitPos := position % SubvectorBits

	return b[subvectorPos].Access(uint8(bitPos))
}

func (b Vector) Ones() uint64 {
	var sum uint64 = 0

	for _, v := range b {
		sum += uint64(bits.OnesCount64(uint64(v)))
	}

	return uint64(sum)
}

func (b Vector) Subvector(position, length uint64) Subvector {
	if length > SubvectorBits {
		panic(fmt.Sprintf("length cant be longer than %d", SubvectorBits))
	}

	subvectorPos := position / SubvectorBits
	bitPos := position % SubvectorBits

	var sub Subvector

	sub = (b[subvectorPos] >> bitPos) & ^(SubvectorMax << length)

	// add overlap
	if bitPos+length > SubvectorBits {
		sub |= b[subvectorPos+1] & ^(SubvectorMax << ((bitPos + length) % SubvectorBits))
	}

	return sub
}

func (b Vector) Bits() uint64 {
	return uint64(len(b) * int(SubvectorBits))
}
