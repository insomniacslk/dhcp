package dhcpv6

import (
	"errors"
	"time"
)

// Response represents a value which Future resolves to
type Response interface {
	Value() DHCPv6
	Error() error
}

// Future is a result of an asynchronous DHCPv6 call
type Future (<-chan Response)

// SuccessFun can be used as a success callback
type SuccessFun func(val DHCPv6) Future

// FailureFun can be used as a failure callback
type FailureFun func(err error) Future

type response struct {
	val DHCPv6
	err error
}

func (r *response) Value() DHCPv6 {
	return r.val
}

func (r *response) Error() error {
	return r.err
}

// NewFuture creates a new future, which can be written to
func NewFuture() chan Response {
	return make(chan Response, 1)
}

// NewResponse creates a new future response
func NewResponse(val DHCPv6, err error) Response {
	return &response{val: val, err: err}
}

// NewSuccessFuture creates a future that resolves to a value
func NewSuccessFuture(val DHCPv6) Future {
	f := NewFuture()
	go func() {
		f <- NewResponse(val, nil)
	}()
	return f
}

// NewFailureFuture creates a future that resolves to an error
func NewFailureFuture(err error) Future {
	f := NewFuture()
	go func() {
		f <- NewResponse(nil, err)
	}()
	return f
}

// Then allows to chain the futures executing appropriate function depending
// on the previous future value
func (f Future) Then(success SuccessFun, failure FailureFun) Future {
	g := NewFuture()
	go func() {
		r := <-f
		if r.Error() != nil {
			r = <-failure(r.Error())
			g <- r
		} else {
			r = <-success(r.Value())
			g <- r
		}
	}()
	return g
}

// OnSuccess allows to chain the futures executing next one only if the first
// one succeeds
func (f Future) OnSuccess(success SuccessFun) Future {
	return f.Then(success, func(err error) Future {
		return NewFailureFuture(err)
	})
}

// OnFailure allows to chain the futures executing next one only if the first
// one fails
func (f Future) OnFailure(failure FailureFun) Future {
	return f.Then(func(val DHCPv6) Future {
		return NewSuccessFuture(val)
	}, failure)
}

// Wait blocks the execution until a future resolves
func (f Future) Wait() (DHCPv6, error) {
	r := <-f
	return r.Value(), r.Error()
}

// WaitTimeout blocks the execution until a future resolves or times out
func (f Future) WaitTimeout(timeout time.Duration) (DHCPv6, error) {
	select {
	case r := <-f:
		return r.Value(), r.Error()
	case <-time.After(timeout):
		return nil, errors.New("Timed out")
	}
}
