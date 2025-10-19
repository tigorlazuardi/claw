package otel

import "runtime"

func CurrentCaller() uintptr {
	return GetCaller(3)
}

func ParentCaller() uintptr {
	return GetCaller(4)
}

func GetCaller(skip int) uintptr {
	pc := make([]uintptr, 1)
	runtime.Callers(skip, pc)
	return pc[0]
}
