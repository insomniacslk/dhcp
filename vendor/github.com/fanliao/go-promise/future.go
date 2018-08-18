/*
Package promise provides a complete promise and future implementation.
A quick start sample:


fu := Start(func()(resp interface{}, err error){
    resp, err = http.Get("http://example.com/")
    return
})
//do somthing...
resp, err := fu.Get()
*/
package promise

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"
)

type callbackType int

const (
	CALLBACK_DONE callbackType = iota
	CALLBACK_FAIL
	CALLBACK_ALWAYS
	CALLBACK_CANCEL
)

//pipe presents a promise that will be chain call
type pipe struct {
	pipeDoneTask, pipeFailTask func(v interface{}) *Future
	pipePromise                *Promise
}

//getPipe returns piped Future task function and pipe Promise by the status of current Promise.
func (this *pipe) getPipe(isResolved bool) (func(v interface{}) *Future, *Promise) {
	if isResolved {
		return this.pipeDoneTask, this.pipePromise
	} else {
		return this.pipeFailTask, this.pipePromise
	}
}

//Canceller is used to check if the future is cancelled
//It usually be passed to the future task function
//for future task function can check if the future is cancelled.
type Canceller interface {
	IsCancelled() bool
	Cancel()
}

//canceller provides an implement of Canceller interface.
//It will be passed to future task function as paramter
type canceller struct {
	f *Future
}

//Cancel sets Future task to CANCELLED status
func (this *canceller) Cancel() {
	this.f.Cancel()
}

//IsCancelled returns true if Future task is cancelld, otherwise false.
func (this *canceller) IsCancelled() (r bool) {
	return this.f.IsCancelled()
}

//futureVal stores the internal state of Future.
type futureVal struct {
	dones, fails, always []func(v interface{})
	cancels              []func()
	pipes                []*pipe
	r                    *PromiseResult
}

//Future provides a read-only view of promise,
//the value is set by using Resolve, Reject and Cancel methods of related Promise
type Future struct {
	Id    int //Id can be used as identity of Future
	final chan struct{}
	//val point to futureVal that stores status of future
	//if need to change the status of future, must copy a new futureVal and modify it,
	//then use CAS to put the pointer of new futureVal
	val unsafe.Pointer
}

//Canceller returns a canceller object related to future.
func (this *Future) Canceller() Canceller {
	return &canceller{this}
}

//IsCancelled returns true if the promise is cancelled, otherwise false
func (this *Future) IsCancelled() bool {
	val := this.loadVal()

	if val != nil && val.r != nil && val.r.Typ == RESULT_CANCELLED {
		return true
	} else {
		return false
	}
}

//SetTimeout sets the future task will be cancelled
//if future is not complete before time out
func (this *Future) SetTimeout(mm int) *Future {
	if mm == 0 {
		mm = 10
	} else {
		mm = mm * 1000 * 1000
	}

	go func() {
		<-time.After((time.Duration)(mm) * time.Nanosecond)
		this.Cancel()
	}()
	return this
}

//GetChan returns a channel than can be used to receive result of Promise
func (this *Future) GetChan() <-chan *PromiseResult {
	c := make(chan *PromiseResult, 1)
	this.OnComplete(func(v interface{}) {
		c <- this.loadResult()
	}).OnCancel(func() {
		c <- this.loadResult()
	})
	return c
}

//Get will block current goroutines until the Future is resolved/rejected/cancelled.
//If Future is resolved, value and nil will be returned
//If Future is rejected, nil and error will be returned.
//If Future is cancelled, nil and CANCELLED error will be returned.
func (this *Future) Get() (val interface{}, err error) {
	<-this.final
	return getFutureReturnVal(this.loadResult())
}

//GetOrTimeout is similar to Get(), but GetOrTimeout will not block after timeout.
//If GetOrTimeout returns with a timeout, timeout value will be true in return values.
//The unit of paramter is millisecond.
func (this *Future) GetOrTimeout(mm uint) (val interface{}, err error, timout bool) {
	if mm == 0 {
		mm = 10
	} else {
		mm = mm * 1000 * 1000
	}

	select {
	case <-time.After((time.Duration)(mm) * time.Nanosecond):
		return nil, nil, true
	case <-this.final:
		r, err := getFutureReturnVal(this.loadResult())
		return r, err, false
	}
}

//Cancel sets the status of promise to RESULT_CANCELLED.
//If promise is cancelled, Get() will return nil and CANCELLED error.
//All callback functions will be not called if Promise is cancalled.
func (this *Future) Cancel() (e error) {
	return this.setResult(&PromiseResult{CANCELLED, RESULT_CANCELLED})
}

//OnSuccess registers a callback function that will be called when Promise is resolved.
//If promise is already resolved, the callback will immediately called.
//The value of Promise will be paramter of Done callback function.
func (this *Future) OnSuccess(callback func(v interface{})) *Future {
	this.addCallback(callback, CALLBACK_DONE)
	return this
}

//OnFailure registers a callback function that will be called when Promise is rejected.
//If promise is already rejected, the callback will immediately called.
//The error of Promise will be paramter of Fail callback function.
func (this *Future) OnFailure(callback func(v interface{})) *Future {
	this.addCallback(callback, CALLBACK_FAIL)
	return this
}

//OnComplete register a callback function that will be called when Promise is rejected or resolved.
//If promise is already rejected or resolved, the callback will immediately called.
//According to the status of Promise, value or error will be paramter of Always callback function.
//Value is the paramter if Promise is resolved, or error is the paramter if Promise is rejected.
//Always callback will be not called if Promise be called.
func (this *Future) OnComplete(callback func(v interface{})) *Future {
	this.addCallback(callback, CALLBACK_ALWAYS)
	return this
}

//OnCancel registers a callback function that will be called when Promise is cancelled.
//If promise is already cancelled, the callback will immediately called.
func (this *Future) OnCancel(callback func()) *Future {
	this.addCallback(callback, CALLBACK_CANCEL)
	return this
}

//Pipe registers one or two functions that returns a Future, and returns a proxy of pipeline Future.
//First function will be called when Future is resolved, the returned Future will be as pipeline Future.
//Secondary function will be called when Futrue is rejected, the returned Future will be as pipeline Future.
func (this *Future) Pipe(callbacks ...interface{}) (result *Future, ok bool) {
	if len(callbacks) == 0 ||
		(len(callbacks) == 1 && callbacks[0] == nil) ||
		(len(callbacks) > 1 && callbacks[0] == nil && callbacks[1] == nil) {
		result = this
		return
	}

	//ensure all callback functions match the spec "func(v interface{}) *Future"
	cs := make([]func(v interface{}) *Future, len(callbacks), len(callbacks))
	for i, callback := range callbacks {
		if c, ok1 := callback.(func(v interface{}) *Future); ok1 {
			cs[i] = c
		} else if c, ok1 := callback.(func() *Future); ok1 {
			cs[i] = func(v interface{}) *Future {
				return c()
			}
		} else if c, ok1 := callback.(func(v interface{})); ok1 {
			cs[i] = func(v interface{}) *Future {
				return Start(func() {
					c(v)
				})
			}
		} else if c, ok1 := callback.(func(v interface{}) (r interface{}, err error)); ok1 {
			cs[i] = func(v interface{}) *Future {
				return Start(func() (r interface{}, err error) {
					r, err = c(v)
					return
				})
			}
		} else if c, ok1 := callback.(func()); ok1 {
			cs[i] = func(v interface{}) *Future {
				return Start(func() {
					c()
				})
			}
		} else if c, ok1 := callback.(func() (r interface{}, err error)); ok1 {
			cs[i] = func(v interface{}) *Future {
				return Start(func() (r interface{}, err error) {
					r, err = c()
					return
				})
			}
		} else {
			ok = false
			return
		}
	}

	for {
		v := this.loadVal()
		r := v.r
		if r != nil {
			result = this
			if r.Typ == RESULT_SUCCESS && cs[0] != nil {
				result = (cs[0](r.Result))
			} else if r.Typ == RESULT_FAILURE && len(cs) > 1 && cs[1] != nil {
				result = (cs[1](r.Result))
			}
		} else {
			newPipe := &pipe{}
			newPipe.pipeDoneTask = cs[0]
			if len(cs) > 1 {
				newPipe.pipeFailTask = cs[1]
			}
			newPipe.pipePromise = NewPromise()

			newVal := *v
			newVal.pipes = append(newVal.pipes, newPipe)

			//use CAS to ensure that the state of Future is not changed,
			//if the state is changed, will retry CAS operation.
			if atomic.CompareAndSwapPointer(&this.val, unsafe.Pointer(v), unsafe.Pointer(&newVal)) {
				result = newPipe.pipePromise.Future
				break
			}
		}
	}
	ok = true

	return
}

//result uses Atomic load to return result of the Future
func (this *Future) loadResult() *PromiseResult {
	val := this.loadVal()
	return val.r
}

//val uses Atomic load to return state value of the Future
func (this *Future) loadVal() *futureVal {
	r := atomic.LoadPointer(&this.val)
	return (*futureVal)(r)
}

//setResult sets the value and final status of Promise, it will only be executed for once
func (this *Future) setResult(r *PromiseResult) (e error) { //r *PromiseResult) {
	defer func() {
		if err := getError(recover()); err != nil {
			e = err
			fmt.Println("\nerror in setResult():", err)
		}
	}()

	e = errors.New("Cannot resolve/reject/cancel more than once")

	for {
		v := this.loadVal()
		if v.r != nil {
			return
		}
		newVal := *v
		newVal.r = r

		//Use CAS operation to ensure that the state of Promise isn't changed.
		//If the state is changed, must get latest state and try to call CAS again.
		//No ABA issue in this case because address of all objects are different.
		if atomic.CompareAndSwapPointer(&this.val, unsafe.Pointer(v), unsafe.Pointer(&newVal)) {
			//Close chEnd then all Get() and GetOrTimeout() will be unblocked
			close(this.final)

			//call callback functions and start the Promise pipeline
			if len(v.dones) > 0 || len(v.fails) > 0 || len(v.always) > 0 || len(v.cancels) > 0 {
				go func() {
					execCallback(r, v.dones, v.fails, v.always, v.cancels)
				}()
			}

			//start the pipeline
			if len(v.pipes) > 0 {
				go func() {
					for _, pipe := range v.pipes {
						pipeTask, pipePromise := pipe.getPipe(r.Typ == RESULT_SUCCESS)
						startPipe(r, pipeTask, pipePromise)
					}
				}()
			}
			e = nil
			break
		}
	}
	return
}

//handleOneCallback registers a callback function
func (this *Future) addCallback(callback interface{}, t callbackType) {
	if callback == nil {
		return
	}
	if (t == CALLBACK_DONE) ||
		(t == CALLBACK_FAIL) ||
		(t == CALLBACK_ALWAYS) {
		if _, ok := callback.(func(v interface{})); !ok {
			panic(errors.New("Callback function spec must be func(v interface{})"))
		}
	} else if t == CALLBACK_CANCEL {
		if _, ok := callback.(func()); !ok {
			panic(errors.New("Callback function spec must be func()"))
		}
	}

	for {
		v := this.loadVal()
		r := v.r
		if r == nil {
			newVal := *v
			switch t {
			case CALLBACK_DONE:
				newVal.dones = append(newVal.dones, callback.(func(v interface{})))
			case CALLBACK_FAIL:
				newVal.fails = append(newVal.fails, callback.(func(v interface{})))
			case CALLBACK_ALWAYS:
				newVal.always = append(newVal.always, callback.(func(v interface{})))
			case CALLBACK_CANCEL:
				newVal.cancels = append(newVal.cancels, callback.(func()))
			}

			//use CAS to ensure that the state of Future is not changed,
			//if the state is changed, will retry CAS operation.
			if atomic.CompareAndSwapPointer(&this.val, unsafe.Pointer(v), unsafe.Pointer(&newVal)) {
				break
			}
		} else {
			if (t == CALLBACK_DONE && r.Typ == RESULT_SUCCESS) ||
				(t == CALLBACK_FAIL && r.Typ == RESULT_FAILURE) ||
				(t == CALLBACK_ALWAYS && r.Typ != RESULT_CANCELLED) {
				callbackFunc := callback.(func(v interface{}))
				callbackFunc(r.Result)
			} else if t == CALLBACK_CANCEL && r.Typ == RESULT_CANCELLED {
				callbackFunc := callback.(func())
				callbackFunc()
			}
			break
		}
	}
}
