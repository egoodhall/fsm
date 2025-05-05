package fsm

type TaskID int64

type State string

type Transition interface {
	To() string
	Data() []byte
}
