package fsm

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/egoodhall/fsm/gen/sqlc"
)

var _ FSM[any, any] = new(fsm[any, any])

type fsm[IN any, OUT any] struct {
	lock sync.Mutex

	// Configurable fields
	name   string
	logger *slog.Logger
	db     sqlc.Querier

	// State machine states
	first FirstTransition[IN, OUT]
	rest  map[State]NthTransition[IN, OUT]

	// Generated / non-configurable fields
	id int64

	// post-start
	tasks chan sqlc.Task
}

func (f *fsm[IN, OUT]) Start(ctx context.Context) (SubmitFunc[IN], context.CancelFunc, error) {
	if !f.lock.TryLock() {
		return nil, nil, errors.New("fsm already started")
	}

	f.tasks = make(chan sqlc.Task, 128)

	go f.process(ctx)

	if err := f.resumeTasks(ctx); err != nil {
		return nil, nil, err
	}

	return f.submit, f.shutdown, nil
}

func (f *fsm[IN, OUT]) shutdown() {
	if f.lock.TryLock() {
		f.lock.Unlock()
		return
	}

	f.logger.Info("Shutting down task processor")
	close(f.tasks)
	f.lock.Unlock()
}

func (f *fsm[IN, OUT]) submit(ctx context.Context, id int64, event IN) (int64, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(event); err != nil {
		return 0, err
	}

	var task sqlc.Task
	var err error
	if id == 0 {
		task, err = f.db.CreateTask(ctx, f.id, buf.Bytes())
	} else {
		task, err = f.db.CreateTaskWithID(ctx, f.id, id, buf.Bytes())
	}

	if err != nil {
		return 0, fmt.Errorf("persist task: %w", err)
	}

	f.enqueue(task)
	return task.ID, nil
}

func (f *fsm[IN, OUT]) process(ctx context.Context) {
	f.logger.Info("Starting task processor")
	for task := range f.tasks {
		if err := f.transition(ctx, task); err != nil {
			f.logger.Error("Failed to transition task to next state", "id", task.ID, "error", err)
		}
	}
}

func (f *fsm[IN, OUT]) resumeTasks(ctx context.Context) error {
	f.logger.Info("Resuming tasks")

	rows, err := f.db.ListTasks(ctx)
	if err != nil {
		return err
	}

	for _, row := range rows {
		state, err := f.db.GetTaskState(ctx, row.ID)
		if errors.Is(err, sql.ErrNoRows) {
			state = "__initial__"
		} else if err != nil {
			return err
		}

		f.logger.Info("Checking task", "id", row.ID, "state", state)
		if State(state) == StateDone {
			continue
		}

		f.logger.Info("Resuming task", "id", row.ID)
		f.enqueue(sqlc.Task(row))

	}

	f.logger.Info("All tasks resumed")
	return nil
}

func (f *fsm[IN, OUT]) enqueue(task sqlc.Task) {
	f.tasks <- task
	f.logger.Info("Enqueued task", "id", task.ID)
}

func (f *fsm[IN, OUT]) transition(ctx context.Context, task sqlc.Task) error {
	transition, err := f.db.GetLastValidTransition(ctx, task.ID)
	if errors.Is(err, sql.ErrNoRows) {
		transition = sqlc.StateTransition{ToState: "__initial__"}
	} else if err != nil {
		return fmt.Errorf("get last valid transition: %w", err)
	}

	var event IN
	if err := gob.NewDecoder(bytes.NewReader(task.Event)).Decode(&event); err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	var (
		out Output
	)
	switch State(transition.ToState) {
	case StateDone:
		f.logger.Info("Completed task", "id", task.ID)
		return nil
	case StateInitial:
		out, err = f.first(ctx, &firstInput[IN]{id: task.ID, event: event, logger: f.logger})
	default:
		var prevout OUT
		if transition.Output != nil {
			if err := gob.NewDecoder(bytes.NewReader(transition.Output)).Decode(&prevout); err != nil {
				return fmt.Errorf("unmarshal output: %w", err)
			}
		}
		if transition, ok := f.rest[State(transition.ToState)]; ok {
			out, err = transition(ctx, &nthInput[IN, OUT]{id: task.ID, event: event, output: prevout, logger: f.logger})
		}
	}
	if err != nil {
		return fmt.Errorf("transition error: %w", err)
	}

	f.logger.Info("Transitioning task", "id", task.ID, "from", transition.ToState, "to", out.NextState())
	if err := f.db.RecordTransition(ctx, task.ID, transition.ToState, string(out.NextState()), out.Data()); err != nil {
		return fmt.Errorf("record transition: %w", err)
	}

	if out.NextState() != StateDone {
		f.enqueue(task)
	}
	return nil
}

var _ FirstInput[any] = new(firstInput[any])

type firstInput[IN any] struct {
	id     int64
	event  IN
	logger *slog.Logger
}

func (r *firstInput[IN]) ID() int64 {
	return r.id
}

func (r *firstInput[IN]) Input() IN {
	return r.event
}

func (r *firstInput[IN]) Logger() *slog.Logger {
	return r.logger
}

var _ NthInput[any, any] = new(nthInput[any, any])

type nthInput[IN any, OUT any] struct {
	id     int64
	event  IN
	output OUT
	logger *slog.Logger
}

func (r *nthInput[IN, OUT]) ID() int64 {
	return r.id
}

func (r *nthInput[IN, OUT]) Input() IN {
	return r.event
}

func (r *nthInput[IN, OUT]) Previous() OUT {
	return r.output
}

func (r *nthInput[IN, OUT]) Logger() *slog.Logger {
	return r.logger
}

func GoTo[T any](state State, data T) (Output, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil && err.Error() == "gob: cannot encode nil value" {
		return &output{state: state, data: buf.Bytes()}, nil
	} else if err != nil {
		return nil, fmt.Errorf("marshal output: %w", err)
	}
	return &output{state: state, data: buf.Bytes()}, nil
}

func Done() (Output, error) {
	return GoTo[any](StateDone, nil)
}

func Error[T any](err error) (Output, error) {
	return &output{state: StateError, data: []byte(err.Error())}, nil
}

var _ Output = new(output)

type output struct {
	state State
	data  []byte
}

func (r *output) NextState() State {
	return r.state
}

func (r *output) Data() []byte {
	return r.data
}
