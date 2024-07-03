package bit

var _ Selectable = (*SelectableBaseline)(nil)

type SelectableBaseline struct {
	Vector
}

// Select implements Selectable.
func (s *SelectableBaseline) Select(alpha bool, n uint64) uint64 {
	var count uint64 = 0

	if n == 0 {
		panic("n hast to be bigger than 0")
	}

	for i := uint64(0); i < s.Bits(); i++ {
		if s.Access(i) == alpha {
			count++
		}

		if count == n {
			return i
		}
	}

	panic("position not found")
}
