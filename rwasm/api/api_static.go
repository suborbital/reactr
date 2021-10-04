package api

import (
	"os"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/rwasm/runtime"
)

func GetStaticFileHandler() runtime.HostFn {
	fn := func(args ...interface{}) (interface{}, error) {
		namePointer := args[0].(int32)
		nameeSize := args[1].(int32)
		ident := args[2].(int32)

		ret := get_static_file(namePointer, nameeSize, ident)

		return ret, nil
	}

	return runtime.NewHostFn("get_static_file", 3, true, fn)
}

func get_static_file(namePtr int32, nameSize int32, ident int32) int32 {
	inst, err := runtime.InstanceForIdentifier(ident, true)
	if err != nil {
		runtime.InternalLogger().Error(errors.Wrap(err, "[rwasm] alert: invalid identifier used, potential malicious activity"))
		return -1
	}

	name := inst.ReadMemory(namePtr, nameSize)

	file, err := inst.Ctx().FileSource.GetStatic(string(name))
	if err != nil {
		runtime.InternalLogger().Error(errors.Wrap(err, "[rwasm] failed to GetStatic"))

		if err == rcap.ErrFileFuncNotSet {
			return -2
		} else if err == os.ErrNotExist {
			return -3
		}

		return -4
	}

	inst.SetFFIResult(file)

	return int32(len(file))
}
