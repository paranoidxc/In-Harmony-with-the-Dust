package syntax

type ComputedToken struct {
	// Offset is the token's start position,
	// defined relative to the computation's start position.
	Offset uint64
	Length uint64
	Role   TokenRole
}
