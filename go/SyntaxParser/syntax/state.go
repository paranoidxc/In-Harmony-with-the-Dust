package syntax

type State interface {
	// Equals returns whether two states are equal.
	//Equals(other State) bool
}

type EmptyState struct{}
