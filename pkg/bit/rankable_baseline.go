package bit

var _ RankableWithSize = (*RankableBaseline)(nil)

type RankableBaseline struct {
	Vector Vector
}

// Size implements RankableWithSize.
func (r *RankableBaseline) Size() uint64 {
	return 0
}

// Rank implements Rankable.
func (r *RankableBaseline) Rank(alpha bool, position uint64) uint64 {
	var rank uint64 = 0

	for i := 0; i < int(position); i++ {
		if r.Vector.Access(uint64(i)) {
			rank++
		}
	}

	if alpha {
		return rank
	}

	return position - rank
}
