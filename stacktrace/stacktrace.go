package stacktrace

import (
	"runtime"
	"strconv"
)

type Frame struct {
	File     string
	Line     int
	Function string
}

func LogFormat(frames []Frame) []string {
	formattedStacks := make([]string, len(frames))
	for i, frame := range frames {
		formattedStacks[i] = frame.File + "(" + frame.Function + ":" + strconv.Itoa(frame.Line) + ")"
	}
	return formattedStacks
}

func AppendStacktrace(frames []Frame, skip, n int) []Frame {
	if n == 0 {
		return frames
	}
	var pc []uintptr
	if n > 0 {
		pc = make([]uintptr, n)
		pc = pc[:runtime.Callers(skip+1, pc)]
	} else {
		n := 0
		pc = make([]uintptr, 10)
		for {
			n += runtime.Callers(skip+n+1, pc[n:])
			if n < len(pc) {
				pc = pc[:n]
				break
			}
			pc = append(pc, 0)
		}
	}
	return AppendCallerFrames(frames, pc, n)
}

func AppendCallerFrames(frames []Frame, callers []uintptr, n int) []Frame {
	if len(callers) == 0 {
		return frames
	}
	runtimeFrames := runtime.CallersFrames(callers)
	for i := 0; n < 0 || i < n; i++ {
		runtimeFrame, more := runtimeFrames.Next()
		frames = append(frames, RuntimeFrame(runtimeFrame))
		if !more {
			break
		}
	}
	return frames
}

func RuntimeFrame(in runtime.Frame) Frame {
	return Frame{
		File:     in.File,
		Function: in.Function,
		Line:     in.Line,
	}
}
