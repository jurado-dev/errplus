package errPlus

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

type errorPlus struct {
	Cause string
	Trace string
	Requeue bool
}

func TraceThis() string{

	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}

func WrapRequeue(input interface{}, trace string) error {

	plus := errorPlus{
		Trace:trace,
		Requeue:true,
	}

	if errorStr, ok := input.(string); ok {
		plus.Cause = errorStr
	}

	if instance, ok := input.(error); ok {
		plus.Cause = instance.Error()
	}

	encoded, err := json.Marshal(plus)
	if err != nil {
		return err
	}

	return errors.New(string(encoded))
}

func Wrap(input interface{},trace string) error {
	plus := errorPlus{
		Trace:trace,
	}

	if errorStr, ok := input.(string); ok {
		plus.Cause = errorStr
	}

	if instance, ok := input.(error); ok {
		plus.Cause = instance.Error()
	}

	encoded, err := json.Marshal(plus)
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

func Decode (input error) (*errorPlus, error) {

	var plus errorPlus

	err := json.Unmarshal([]byte(input.Error()), &plus)
	if err != nil {
		return nil, err
	}

	return &plus, nil
}

