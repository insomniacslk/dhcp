package promise

import (
	"math/rand"
	"unsafe"
)

var (
	CANCELLED error = &CancelledError{}
)

//CancelledError present the Future object is cancelled.
type CancelledError struct {
}

func (e *CancelledError) Error() string {
	return "Task be cancelled"
}

//resultType present the type of Future final status.
type resultType int

const (
	RESULT_SUCCESS resultType = iota
	RESULT_FAILURE
	RESULT_CANCELLED
)

//PromiseResult presents the result of a promise.
//If Typ is RESULT_SUCCESS, Result field will present the returned value of Future task.
//If Typ is RESULT_FAILURE, Result field will present a related error .
//If Typ is RESULT_CANCELLED, Result field will be null.
type PromiseResult struct {
	Result interface{} //result of the Promise
	Typ    resultType  //success, failure, or cancelled?
}

//Promise presents an object that acts as a proxy for a result.
//that is initially unknown, usually because the computation of its
//value is yet incomplete (refer to wikipedia).
//You can use Resolve/Reject/Cancel to set the final result of Promise.
//Future can return a read-only placeholder view of result.
type Promise struct {
	*Future
}

//Cancel sets the status of promise to RESULT_CANCELLED.
//If promise is cancelled, Get() will return nil and CANCELLED error.
//All callback functions will be not called if Promise is cancalled.
func (this *Promise) Cancel() (e error) {
	return this.Future.Cancel()
}

//Resolve sets the value for promise, and the status will be changed to RESULT_SUCCESS.
//if promise is resolved, Get() will return the value and nil error.
func (this *Promise) Resolve(v interface{}) (e error) {
	return this.setResult(&PromiseResult{v, RESULT_SUCCESS})
}

//Resolve sets the error for promise, and the status will be changed to RESULT_FAILURE.
//if promise is rejected, Get() will return nil and the related error value.
func (this *Promise) Reject(err error) (e error) {
	return this.setResult(&PromiseResult{err, RESULT_FAILURE})
}

//OnSuccess registers a callback function that will be called when Promise is resolved.
//If promise is already resolved, the callback will immediately called.
//The value of Promise will be paramter of Done callback function.
func (this *Promise) OnSuccess(callback func(v interface{})) *Promise {
	this.Future.OnSuccess(callback)
	return this
}

//OnFailure registers a callback function that will be called when Promise is rejected.
//If promise is already rejected, the callback will immediately called.
//The error of Promise will be paramter of Fail callback function.
func (this *Promise) OnFailure(callback func(v interface{})) *Promise {
	this.Future.OnFailure(callback)
	return this
}

//OnComplete register a callback function that will be called when Promise is rejected or resolved.
//If promise is already rejected or resolved, the callback will immediately called.
//According to the status of Promise, value or error will be paramter of Always callback function.
//Value is the paramter if Promise is resolved, or error is the paramter if Promise is rejected.
//Always callback will be not called if Promise be called.
func (this *Promise) OnComplete(callback func(v interface{})) *Promise {
	this.Future.OnComplete(callback)
	return this
}

//OnCancel registers a callback function that will be called when Promise is cancelled.
//If promise is already cancelled, the callback will immediately called.
func (this *Promise) OnCancel(callback func()) *Promise {
	this.Future.OnCancel(callback)
	return this
}

//NewPromise is factory function for Promise
func NewPromise() *Promise {
	val := &futureVal{
		make([]func(v interface{}), 0, 8),
		make([]func(v interface{}), 0, 8),
		make([]func(v interface{}), 0, 4),
		make([]func(), 0, 2),
		make([]*pipe, 0, 4), nil,
	}
	f := &Promise{
		&Future{
			rand.Int(),
			make(chan struct{}),
			unsafe.Pointer(val),
		},
	}
	return f
}
