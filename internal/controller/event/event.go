package event

type Event interface {
	Type() int
	String() string
}
