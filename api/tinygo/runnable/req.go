package runnable

// #include "reactr.h"
import "C"

type FieldType int32

const (
	FieldMeta FieldType = iota
	FieldBody
	FieldHeader
	FieldParams
	FieldState
)

func getField(fieldType FieldType, key string) []byte {
	ptr, size := rawSlicePointer([]byte(key))

	res, err := result(C.request_get_field(int32(fieldType), ptr, size, ident()))
	if err != nil {
		return []byte{}
	}
	return res
}

func setField(fieldType FieldType, key string, value string) ([]byte, error) {
	keyPtr, keySize := rawSlicePointer([]byte(key))
	valPtr, valSize := rawSlicePointer([]byte(value))

	return result(C.request_set_field(int32(fieldType), keyPtr, keySize, valPtr, valSize, ident()))
}

func Method() string {
	return string(getField(FieldMeta, "method"))
}

func SetMethod(value string) error {
	_, err := setField(FieldMeta, "method", value)
	return err
}

func URL() string {
	return string(getField(FieldMeta, "url"))
}

func SetURL(value string) error {
	_, err := setField(FieldMeta, "url", value)
	return err
}

func ID() string {
	return string(getField(FieldMeta, "id"))
}

func Body() []byte {
	return getField(FieldMeta, "body")
}

func BodyString() string {
	return string(Body())
}

func SetBody(value string) error {
	_, err := setField(FieldBody, "body", value)
	return err
}

func BodyField(key string) string {
	return string(getField(FieldBody, key))
}

func SetBodyField(key, value string) error {
	_, err := setField(FieldBody, key, value)
	return err
}

func Header(key string) string {
	return string(getField(FieldHeader, key))
}

func SetHeader(key, value string) error {
	_, err := setField(FieldHeader, key, value)
	return err
}

func URLParam(key string) string {
	return string(getField(FieldParams, key))
}

func SetURLParam(key, value string) error {
	_, err := setField(FieldParams, key, value)
	return err
}

func StateString(key string) string {
	return string(State(key))
}

func State(key string) []byte {
	return getField(FieldState, key)
}

func SetState(key, value string) error {
	_, err := setField(FieldState, key, value)
	return err
}
