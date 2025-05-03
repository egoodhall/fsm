package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/egoodhall/fsm"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	machine, err := fsm.New[string, string]("example").
		InitialState(func(ctx context.Context, req fsm.FirstInput[string]) (fsm.Output, error) {
			if req.Input() != "a" {
				return fsm.Done()
			}
			return fsm.GoTo("a", "_")
		}).
		AddState("a", func(ctx context.Context, req fsm.NthInput[string, string]) (fsm.Output, error) {
			if req.Previous() == "_" {
				return fsm.GoTo("b", req.Previous()+"a")
			}
			return fsm.Done()
		}).
		AddState("b", func(ctx context.Context, req fsm.NthInput[string, string]) (fsm.Output, error) {
			return fsm.GoTo("a", req.Previous()+"b")
		}).
		Build(ctx, fsm.Logger[string, string](slog.Default()), fsm.DB[string, string]("fsm.db"))

	if err != nil {

		slog.Error("Failed to build machine", "error", err)

	}

	slog.Info("Starting machine")
	submit, stop, err := machine.Start(ctx)
	if err != nil {
		slog.Error("Failed to start machine", "error", err)
	}
	defer stop()

	if err := submit(ctx, "1", "a"); err != nil {
		slog.Error("Failed to submit task", "error", err)
	}

	if err := submit(ctx, "2", "a"); err != nil {
		slog.Error("Failed to submit task", "error", err)
	}

	if err := submit(ctx, "3", "b"); err != nil {
		slog.Error("Failed to submit task", "error", err)
	}

	<-ctx.Done()
}
