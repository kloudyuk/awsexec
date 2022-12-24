package awsexec

import (
	"fmt"
	"reflect"
	"sync"
)

type result struct {
	m       sync.Mutex
	results reflect.Value
}

func (r *result) Add(v interface{}) {
	result := reflect.ValueOf(v)
	r.m.Lock()
	if result.Kind().String() == "slice" {
		r.results.Set(reflect.AppendSlice(r.results, result))
	} else {
		r.results.Set(reflect.Append(r.results, result))
	}
	r.m.Unlock()
}

// errs is a custom error type that holds a slice of errors and
// can be returned as a single error as it implements the error interface
type execErr struct {
	m       sync.Mutex
	errList []error
}

func (e *execErr) Error() string {
	var errStr string
	for _, err := range e.errList {
		errStr += fmt.Sprintf("%s\n", err)
	}
	return errStr
}

func (e *execErr) Add(err error) {
	e.m.Lock()
	e.errList = append(e.errList, err)
	e.m.Unlock()
}

func (e *execErr) Len() int {
	return len(e.errList)
}
