package internal

import (
	"context"
	"time"

	"github.com/tigorlazuardi/claw/lib/otel"
	"github.com/urfave/cli/v3"
)

type detachedContext struct {
	context.Context
}

func (de detachedContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (de detachedContext) Done() <-chan struct{} {
	return nil
}

func (de detachedContext) Err() error {
	return nil
}

func (de detachedContext) Value(key any) any {
	return de.Context.Value(key)
}

func After(ctx context.Context, cmd *cli.Command) error {
	ctx, cancel := context.WithTimeout(detachedContext{ctx}, 10*time.Second)
	defer cancel()
	return otel.ShutdownAll(ctx)
}
