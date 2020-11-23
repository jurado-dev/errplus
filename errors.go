package errPlus

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

type errorPlus struct {
	Cause   string
	Trace   string
	Code    int
	Requeue bool
}

func TraceThis() string {

	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}

func WrapRequeue(input interface{}, trace string) error {

	var ePlus errorPlus

	if custom, ok := input.(string); ok {
		ePlus.Cause = custom
	}

	if err, ok := input.(error); ok {
		ePlus.Cause = err.Error()
	}

	if instance, ok := input.(errorPlus); ok {
		ePlus.Code = instance.Code
		ePlus.Cause = instance.Cause
	}

	ePlus.Requeue = true
	ePlus.Trace = trace

	encoded, err := json.Marshal(ePlus)
	if err != nil {
		return err
	}

	return errors.New(string(encoded))
}

func Wrap(input interface{}, trace string) error {

	var ePlus errorPlus

	if custom, ok := input.(string); ok {
		ePlus.Cause = custom
	}

	if err, ok := input.(error); ok {
		ePlus.Cause = err.Error()
	}

	if instance, ok := input.(errorPlus); ok {
		ePlus.Code = instance.Code
		ePlus.Cause = instance.Cause
	}

	ePlus.Trace = trace

	encoded, err := json.Marshal(ePlus)
	if err != nil {
		return err
	}

	return errors.New(string(encoded))
}

func GetCause(input error) (string, error) {

	if input == nil {
		return "", nil
	}

	plus, err := Decode(input)
	if err != nil {
		return "", err
	}

	return plus.Cause, nil
}

func GetString(input error) (string, error) {

	if input == nil {
		return "", nil
	}

	plus, err := Decode(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("trace: %s | requeue: %v | cause: %s", plus.Trace, plus.Requeue, plus.Cause), nil
}

func GetRequeue(input error) (bool, error) {

	if input == nil {
		return false, nil
	}

	plus, err := Decode(input)
	if err != nil {
		return false, err
	}

	return plus.Requeue, nil
}

func GetCode(input error) (int, error) {
	if input == nil {
		return 0, nil
	}
	ePlus, err := Decode(input)
	if err != nil {
		return 0, err
	}
	return ePlus.Code, nil
}

func Decode(input error) (*errorPlus, error) {

	var plus errorPlus

	err := json.Unmarshal([]byte(input.Error()), &plus)
	if err != nil {
		return nil, err
	}

	return &plus, nil
}

func ErrWithCode(cause string, code int) interface{} {
	return errorPlus{Cause: cause, Code: code}
}

