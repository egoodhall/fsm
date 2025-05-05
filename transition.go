package fsm

type State string

type Transition interface {
	To() string
	Data() []byte
}
