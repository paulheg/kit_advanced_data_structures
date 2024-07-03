package bit

import (
	"cmp"
	"math"
)

const InterleavedSubvectorCount uint64 = 7

// Define 512 bit datastructure that will fit in one cache line
type InterleavedVectorLine struct {
	PreSum uint64
	Vec    [InterleavedSubvectorCount]Subvector
}

var _ RankSelectVector = (*InterleavedVector)(nil)
var _ Setable = (*InterleavedVector)(nil)

// Create the interleaved datastructure by copying the vec Vector
// Also create the pre sums
func NewInterleavedVector(vec Vector) *InterleavedVector {
	lines := uint64(math.Ceil(float64(len(vec)) / float64(InterleavedSubvectorCount)))

	intlVec := &InterleavedVector{
		vec: make([]InterleavedVectorLine, lines),
	}

	var totalSum uint64

	for i := range lines {
		startPos := i * InterleavedSubvectorCount
		endPos := startPos + InterleavedSubvectorCount

		if endPos > uint64(len(vec)) {
			endPos = uint64(len(vec))
		}

		lineSlice := vec[startPos:endPos]
		copy(intlVec.vec[i].Vec[:], lineSlice)

		intlVec.vec[i].PreSum = totalSum

		totalSum += Vector(lineSlice).Ones()
	}

	return intlVec
}

func NewInterleavedVectorNoPrecompute(vec Vector) *InterleavedVector {
	lines := uint64(math.Ceil(float64(len(vec)) / float64(InterleavedSubvectorCount)))

	intlVec := &InterleavedVector{
		vec: make([]InterleavedVectorLine, lines),
	}

	for i := range lines {
		startPos := i * InterleavedSubvectorCount
		endPos := startPos + InterleavedSubvectorCount

		if endPos > uint64(len(vec)) {
			endPos = uint64(len(vec))
		}

		lineSlice := vec[startPos:endPos]
		copy(intlVec.vec[i].Vec[:], lineSlice)
	}

	return intlVec
}

type InterleavedVector struct {
	vec []InterleavedVectorLine
}

// Calculate the pre sums on an otherwise filled InterleavedVector
func (i *InterleavedVector) Precompute() {
	var sum uint64

	for j := range len(i.vec) {
		i.vec[j].PreSum = sum
		sum += Vector(i.vec[j].Vec[:]).Ones()
	}
}

// Find the correct subvector and also get its inner position
func (i *InterleavedVector) GetSubvector(position uint64) (innerPos uint8, subvec *Subvector) {
	subvectorPos := position / SubvectorBits
	innerPos = uint8(position % SubvectorBits)

	interleavedVectorLinePos := subvectorPos / InterleavedSubvectorCount
	interleavedSubVectorPos := subvectorPos % InterleavedSubvectorCount

	subvec = &i.vec[interleavedVectorLinePos].Vec[interleavedSubVectorPos]
	return
}

// Set implements Setable.
func (i *InterleavedVector) Set(position uint64) {
	pos, sv := i.GetSubvector(position)
	sv.Set(pos)
}

// Access implements RankSelectVector.
func (i *InterleavedVector) Access(position uint64) bool {
	ipos, sv := i.GetSubvector(position)
	return sv.Access(ipos)
}

// Rank implements RankSelectVector.(number before)
func (i *InterleavedVector) Rank(alpha bool, position uint64) uint64 {
	subvectorPos := position / SubvectorBits
	innerSubVecPos := position % SubvectorBits

	interleavedVectorLinePos := subvectorPos / InterleavedSubvectorCount
	interleavedSubVectorPos := subvectorPos % InterleavedSubvectorCount

	line := i.vec[interleavedVectorLinePos]
	rank := line.PreSum

	if interleavedSubVectorPos > 0 {
		rank += Vector(line.Vec[0:(interleavedSubVectorPos)]).Ones()
	}

	rank += uint64(line.Vec[interleavedSubVectorPos].Rank(true, uint8(innerSubVecPos)))

	if alpha {
		return rank
	} else {
		return position - rank
	}
}

// Select implements RankSelectVector.(nth one)
func (i *InterleavedVector) Select(alpha bool, n uint64) uint64 {

	linePos, _ := i.BinarySearch(alpha, n)

	if linePos > 0 {
		linePos--
	}

	var prevSum uint64

	line := i.vec[linePos]

	// select the previous block if available
	if linePos > 0 {
		prevSum = line.PreSum
		if !alpha {
			prevSum = uint64(linePos*InterleavedSubvectorCount*SubvectorBits) - prevSum
		}
	}

	subvecPos := 0
	for {
		currentCount := uint64(line.Vec[subvecPos].Ones())
		if !alpha {
			currentCount = SubvectorBits - currentCount
		}

		if prevSum+currentCount >= n {
			break
		}
		subvecPos++
		prevSum += currentCount
	}

	// calculate the position inside the block
	var innerPos uint8 = uint8(n - prevSum)
	innerValue := uint64(line.Vec[subvecPos].Select(alpha, innerPos))
	return innerValue + uint64(linePos)*InterleavedSubvectorCount*SubvectorBits + uint64(subvecPos)*SubvectorBits
}

func (c *InterleavedVector) BinarySearch(alpha bool, target uint64) (uint64, bool) {

	// Inlining is faster than calling BinarySearchFunc with a lambda.
	n := uint64(len(c.vec))
	// Define x[-1] < target and x[n] >= target.
	// Invariant: x[i-1] < target, x[j] >= target.
	i, j := uint64(0), n
	for i < j {
		h := uint64(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j

		value := c.vec[h].PreSum
		if !alpha {
			value = (h * InterleavedSubvectorCount * SubvectorBits) - value
		}

		if cmp.Less(value, target) {
			i = h + 1 // preserves x[i-1] < target
		} else {
			j = h // preserves x[j] >= target
		}
	}
	// i == j, x[i-1] < target, and x[j] (= x[i]) >= target  =>  answer is i.
	return i, i < n && c.vec[i].PreSum == target
}

// Overhead implements RankSelectVector.
func (i *InterleavedVector) Overhead() uint64 {
	// 1 uint64 per line
	return uint64(len(i.vec) * 64)
}

// Size implements RankSelectVector.
func (i *InterleavedVector) Size() uint64 {
	// each line is 512 bit in size
	return uint64(len(i.vec) * 512)
}
