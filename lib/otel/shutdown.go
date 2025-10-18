package otel

import "context"

type Shutdown func(ctx context.Context) error

var Shutdowns []Shutdown

func ShutdownAll(ctx context.Context) error {
	for _, shutdown := range Shutdowns {
		if err := shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
