package bit

type RankableWithSize interface {
	Rankable
	Sizable
}

type Rankable interface {
	// Number of alphas before position
	Rank(alpha bool, position uint64) uint64
}
