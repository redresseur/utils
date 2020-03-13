package errors

import (
	"runtime"
	"strings"
)

const (
	maximumCallerDepth int = 25
	minimumCallerDepth int = 3
)

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func callers() []runtime.Frame {
	stacks := []runtime.Frame{}
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for {
		f, again := frames.Next()
		stacks = append(stacks, f)
		if !again {
			break
		}
	}

	return stacks
}
