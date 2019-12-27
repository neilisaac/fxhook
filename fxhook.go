package fxhook

import (
	"context"
	"errors"

	"go.uber.org/fx"
)

func CtxErrFunc(callback func(context.Context) error) fx.Hook {
	runCtx, cancel := context.WithCancel(context.Background())
	result := make(chan error)

	return fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				defer close(result)
				result <- callback(runCtx)
			}()
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			cancel()
			select {
			case err := <-result:
				if errors.Is(err, runCtx.Err()) {
					return nil
				}
				return err
			case <-stopCtx.Done():
				return stopCtx.Err()
			}
		},
	}
}
