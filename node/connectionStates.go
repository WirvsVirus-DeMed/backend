package node

type ConnectionState byte

const (
	Closed      ConnectionState = 0b000
	Established ConnectionState = 0b001
	Ready       ConnectionState = 0b010
)
