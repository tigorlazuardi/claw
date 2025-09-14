package server

import (
	"context"
	"runtime"
)

type enpoindtInfoKey struct{}

type endpointInfo struct {
	pc uintptr
}

func registerEndpointInfo(ctx context.Context) {
	epinfo := getEndpointInfo(ctx)
	if epinfo == nil {
		return
	}
	pcs := make([]uintptr, 1)
	n := runtime.Callers(2, pcs)
	if n == 0 {
		return
	}
	epinfo.pc = pcs[0]
}

func getEndpointInfo(ctx context.Context) *endpointInfo {
	v, _ := ctx.Value(enpoindtInfoKey{}).(*endpointInfo)
	return v
}

func withEndpointInfo(ctx context.Context) context.Context {
	return context.WithValue(ctx, enpoindtInfoKey{}, &endpointInfo{})
}

func getEnpointInfoPC(ctx context.Context) uintptr {
	epinfo := getEndpointInfo(ctx)
	if epinfo == nil {
		return 0
	}
	return epinfo.pc
}
