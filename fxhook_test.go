package fxhook

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/fx"
)

func TestCtxErrFuncError(t *testing.T) {
	started := make(chan struct{})
	finished := make(chan struct{})
	expectedError := errors.New("expected")

	callback := func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		close(finished)
		return expectedError
	}

	register := func(lc fx.Lifecycle) {
		lc.Append(CtxErrFunc(callback))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app := fx.New(fx.Invoke(register))
	if err := app.Start(ctx); err != nil {
		t.Fatal("Start error:", err)
	}

	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("callback did not start")
	}

	if err := app.Stop(ctx); err != expectedError {
		t.Error("Stop error:", err)
	}

	select {
	case <-finished:
	case <-ctx.Done():
		t.Error("callback did not get informed to stop")
	}
}
