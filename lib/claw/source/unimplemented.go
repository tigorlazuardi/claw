package source

import (
	"context"
	"errors"
)

var _ Source = (*UnimplementedSource)(nil)

const UnimplementedSourceName = "unimplemented.source"

var ErrUnimplementedSource = errors.New("source is not implemented")

// UnimplementedSource is a stub implementation of the Source interface.
//
// Embed this struct to have a default implementation for various methods and override
// only the methods you need.
//
// This is also useful to avoid breaking changes when new methods are added to the Source interface.
//
// Example:
//
//	type MySource struct {
//		source.UnimplementedSource // Keep this at the top so all methods are properly overridden
//		// other fields...
//	}
//
//	func (ms *MySource) Name() string {
//		return "my.source.v1"
//	}
//
//	func (ms *MySource) Run(ctx context.Context, request source.Request) (source.Response, error) {
//		// implementation...
//		return source.Response{}, nil
//	}
type UnimplementedSource struct{}

func (UnimplementedSource) HaveScheduleConflictCheck() bool {
	return false
}

func (UnimplementedSource) ScheduleConflictCheck(req ScheduleConflictCheckRequest) string {
	return ""
}

func (UnimplementedSource) Name() string {
	return UnimplementedSourceName
}

func (UnimplementedSource) Run(ctx context.Context, request Request) (Response, error) {
	return Response{}, ErrUnimplementedSource
}

func (UnimplementedSource) Description() string {
	return ""
}

func (UnimplementedSource) DisplayName() string {
	return "Unknown"
}

func (UnimplementedSource) Author() string {
	return ""
}

func (UnimplementedSource) AuthorURL() string {
	return ""
}

func (UnimplementedSource) RequireParameter() bool {
	return false
}

func (UnimplementedSource) ParameterHelp() string {
	return ""
}

func (UnimplementedSource) ParameterPlaceholder() string {
	return ""
}

func (UnimplementedSource) ValidateTransformParameter(ctx context.Context, param string) (transformed string, err error) {
	return param, nil
}

func (UnimplementedSource) DefaultCountback() int {
	return 0
}
