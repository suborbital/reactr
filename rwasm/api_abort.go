package rwasm

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func abortHandler() *HostFn {
	fn := func(args ...wasmer.Value) (interface{}, error) {
		message := args[0].I32()
		fileName := args[1].I32()
		lineNumber := args[2].I32()
		columnNumber := args[3].I32()
		ident := args[4].I32()

		abort(message, fileName, lineNumber, columnNumber, ident)

		return nil, nil
	}

	return newHostFn("abort", 4, false, fn)
}

func abort(message int32, fileName int32, lineNumber int32, columnNumber int32, ident int32) {
	inst, err := instanceForIdentifier(ident, false)
	if err != nil {
		logger.Error(errors.Wrap(err, "[rwasm] alert: invalid identifier used, potential malicious activity"))
		return
	}

	errMsg := fmt.Sprintf("aborted: line: %d, col: %d", lineNumber, columnNumber)

	runErr := rt.RunErr{Code: 1, Message: errMsg}

	inst.errChan <- runErr
}
