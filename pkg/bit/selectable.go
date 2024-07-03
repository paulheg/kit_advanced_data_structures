package bit

type Selectable interface {
	// Position of the n'th alpha
	Select(alpha bool, n uint64) uint64
}

type SelectableWithSize interface {
	Selectable
	Sizable
}
