package api

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rwasm/runtime"
)

func DBExecHandler() runtime.HostFn {
	fn := func(args ...interface{}) (interface{}, error) {
		queryType := args[0].(int32)
		namePointer := args[1].(int32)
		nameSize := args[2].(int32)
		ident := args[3].(int32)

		ret := db_exec(queryType, namePointer, nameSize, ident)

		return ret, nil
	}

	return runtime.NewHostFn("db_exec", 4, true, fn)
}

func db_exec(queryType, namePointer, nameSize, identifier int32) int32 {
	inst, err := runtime.InstanceForIdentifier(identifier, false)
	if err != nil {
		runtime.InternalLogger().Error(errors.Wrap(err, "[rwasm] alert: invalid identifier used, potential malicious activity"))
		return -1
	}

	nameBytes := inst.ReadMemory(namePointer, nameSize)
	name := string(nameBytes)

	_, err = inst.Ctx().Database.ExecInsertQuery(name, nil)
	if err != nil {
		runtime.InternalLogger().ErrorString("[rwasm] failed to ExexInsertQuery", name, err.Error())

		res, _ := inst.Ctx().SetFFIResult(nil, err)
		return res.FFISize()
	}

	res, _ := inst.Ctx().SetFFIResult([]byte("success!"), nil)

	return res.FFISize()
}
