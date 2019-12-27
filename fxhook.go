package fxhook

import (
	"context"

	"go.uber.org/fx"
)

func Go(callback func(context.Context) error) fx.Hook {
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
				return err
			case <-stopCtx.Done():
				return stopCtx.Err()
			}
		},
	}
}
