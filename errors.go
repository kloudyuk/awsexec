package awsexec

import (
	"fmt"
	"sync"
)

// execErr is a custom error type that holds a slice of errors and
// can be returned as a single error as it implements the error interface
type execErr struct {
	sync.Mutex
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
	e.Lock()
	e.errList = append(e.errList, err)
	e.Unlock()
}

func (e *execErr) Len() int {
	return len(e.errList)
}
