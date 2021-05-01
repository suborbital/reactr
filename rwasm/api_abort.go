package rwasm

import (
	"fmt"

	"github.com/wasmerio/wasmer-go/wasmer"
)

func abortHandler() *HostFn {
	fn := func(args ...wasmer.Value) (interface{}, error) {
		message := args[0].I32()
		fileName := args[1].I32()
		lineNumber := args[2].I32()
		columnNumber := args[3].I32()

		abort(message, fileName, lineNumber, columnNumber)

		return nil, nil
	}

	return newHostFn("abort", 4, false, fn)
}

func abort(message int32, fileName int32, lineNumber int32, columnNumber int32) {
	errMsg := fmt.Sprintf("instance aborted: line: %d, col: %d", lineNumber, columnNumber)

	logger.ErrorString(errMsg)
}
