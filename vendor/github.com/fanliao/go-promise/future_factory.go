package promise

import (
	"sync/atomic"
)

type anyPromiseResult struct {
	result interface{}
	i      int
}

//Start start a goroutines to execute task function
//and return a Future that presents the result.
//If option paramter is true, the act function will be sync called.
//Type of act can be any of below four types:
//  func() (r interface{}, err error):
//     if err returned by act != nil or panic error, then Future will be rejected with error,
//     otherwise be resolved with r.
//  func():
//     if act panic error, then Future will be rejected, otherwise be resolved with nil.
//  func(c promise.Canceller) (r interface{}, err error):
//     if err returned by act != nil or panic error,
//     then Future will be rejected with err, otherwise be resolved with r.
//     We can check c.IsCancelled() to decide whether need to exit act function
//  func(promise.Canceller):
//     if act panic error, then Future will be rejected with error, otherwise be resolved with nil.
//     We can check c.IsCancelled() to decide whether need to exit act function.
//  error:
//     Future will be rejected with error immediately
//  other value:
//     Future will be resolved with value immediately
func Start(act interface{}, syncs ...bool) *Future {
	pr := NewPromise()
	if f, ok := act.(*Future); ok {
		return f
	}

	if action := getAct(pr, act); action != nil {
		if syncs != nil && len(syncs) > 0 && !syncs[0] {
			//sync call
			r, err := action()
			if pr.IsCancelled() {
				pr.Cancel()
			} else {
				if err == nil {
					pr.Resolve(r)
				} else {
					pr.Reject(err)
				}
			}
		} else {
			//async call
			go func() {
				r, err := action()
				if pr.IsCancelled() {
					pr.Cancel()
				} else {
					if err == nil {
						pr.Resolve(r)
					} else {
						pr.Reject(err)
					}
				}
			}()
		}
	}

	return pr.Future
}

//Wrap return a Future that presents the wrapped value
func Wrap(value interface{}) *Future {
	pr := NewPromise()
	if e, ok := value.(error); !ok {
		pr.Resolve(value)
	} else {
		pr.Reject(e)
	}

	return pr.Future
}

//WhenAny returns a Future.
//If any Future is resolved, this Future will be resolved and return result of resolved Future.
//Otherwise will rejected with results slice returned by all Futures
//Legit types of act are same with Start function
func WhenAny(acts ...interface{}) *Future {
	return WhenAnyMatched(nil, acts...)
}

//WhenAnyMatched returns a Future.
//If any Future is resolved and match the predicate, this Future will be resolved and return result of resolved Future.
//If all Futures are cancelled, this Future will be cancelled.
//Otherwise will rejected with a NoMatchedError included results slice returned by all Futures
//Legit types of act are same with Start function
func WhenAnyMatched(predicate func(interface{}) bool, acts ...interface{}) *Future {
	if predicate == nil {
		predicate = func(v interface{}) bool { return true }
	}

	fs := make([]*Future, len(acts))
	for i, act := range acts {
		fs[i] = Start(act)
	}

	nf, rs := NewPromise(), make([]interface{}, len(fs))
	if len(acts) == 0 {
		nf.Resolve(nil)
	}

	chFails, chDones := make(chan anyPromiseResult), make(chan anyPromiseResult)

	go func() {
		for i, f := range fs {
			k := i
			f.OnSuccess(func(v interface{}) {
				defer func() { _ = recover() }()
				chDones <- anyPromiseResult{v, k}
			}).OnFailure(func(v interface{}) {
				defer func() { _ = recover() }()
				chFails <- anyPromiseResult{v, k}
			}).OnCancel(func() {
				defer func() { _ = recover() }()
				chFails <- anyPromiseResult{CANCELLED, k}
			})
		}
	}()

	if len(fs) == 1 {
		select {
		case r := <-chFails:
			if _, ok := r.result.(CancelledError); ok {
				nf.Cancel()
			} else {
				nf.Reject(newNoMatchedError1(r.result))
			}
		case r := <-chDones:
			if predicate(r.result) {
				nf.Resolve(r.result)
			} else {
				nf.Reject(newNoMatchedError1(r.result))
			}
		}
	} else {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					nf.Reject(newErrorWithStacks(e))
				}
			}()

			j := 0
			for {
				select {
				case r := <-chFails:
					rs[r.i] = getError(r.result)
				case r := <-chDones:
					if predicate(r.result) {
						//try to cancel other futures
						for _, f := range fs {
							f.Cancel()
						}

						//close the channel for avoid the send side be blocked
						closeChan := func(c chan anyPromiseResult) {
							defer func() { _ = recover() }()
							close(c)
						}
						closeChan(chDones)
						closeChan(chFails)

						//Resolve the future and return result
						nf.Resolve(r.result)
						return
					} else {
						rs[r.i] = r.result
					}
				}

				if j++; j == len(fs) {
					m := 0
					for _, r := range rs {
						switch val := r.(type) {
						case CancelledError:
						default:
							m++
							_ = val
						}
					}
					if m > 0 {
						nf.Reject(newNoMatchedError(rs))
					} else {
						nf.Cancel()
					}
					break
				}
			}
		}()
	}
	return nf.Future
}

//WhenAll receives function slice and returns a Future.
//If all Futures are resolved, this Future will be resolved and return results slice.
//Otherwise will rejected with results slice returned by all Futures
//Legit types of act are same with Start function
func WhenAll(acts ...interface{}) (fu *Future) {
	pr := NewPromise()
	fu = pr.Future

	if len(acts) == 0 {
		pr.Resolve([]interface{}{})
		return
	}

	fs := make([]*Future, len(acts))
	for i, act := range acts {
		fs[i] = Start(act)
	}
	fu = whenAllFuture(fs...)
	return
}

//WhenAll receives Futures slice and returns a Future.
//If all Futures are resolved, this Future will be resolved and return results slice.
//If any Future is cancelled, this Future will be cancelled.
//Otherwise will rejected with results slice returned by all Futures.
//Legit types of act are same with Start function
func whenAllFuture(fs ...*Future) *Future {
	wf := NewPromise()
	rs := make([]interface{}, len(fs))

	if len(fs) == 0 {
		wf.Resolve([]interface{}{})
	} else {
		n := int32(len(fs))
		cancelOthers := func(j int) {
			for k, f1 := range fs {
				if k != j {
					f1.Cancel()
				}
			}
		}

		go func() {
			isCancelled := int32(0)
			for i, f := range fs {
				j := i

				f.OnSuccess(func(v interface{}) {
					rs[j] = v
					if atomic.AddInt32(&n, -1) == 0 {
						wf.Resolve(rs)
					}
				}).OnFailure(func(v interface{}) {
					if atomic.CompareAndSwapInt32(&isCancelled, 0, 1) {
						//try to cancel all futures
						cancelOthers(j)

						//errs := make([]error, 0, 1)
						//errs = append(errs, v.(error))
						e := newAggregateError1("Error appears in WhenAll:", v)
						wf.Reject(e)
					}
				}).OnCancel(func() {
					if atomic.CompareAndSwapInt32(&isCancelled, 0, 1) {
						//try to cancel all futures
						cancelOthers(j)

						wf.Cancel()
					}
				})
			}
		}()
	}

	return wf.Future
}
