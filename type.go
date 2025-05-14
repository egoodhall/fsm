package fsm

type State string

// StateError is used to indicate an error during a state transition.
const StateError State = "__error__"

type TaskID int64
