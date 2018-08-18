[home]: github.com/fanliao/go-promise

go-promise is a Go promise and future library.

Inspired by [Futures and promises]()

## Installation

    $ go get github.com/fanliao/go-promise

## Features

* Future and Promise

  * ```NewPromise()```
  * ```promise.Future```

* Promise and Future callbacks

  * ```.OnSuccess(v interface{})```
  * ```.OnFailure(v interface{})```
  * ```.OnComplete(v interface{})```
  * ```.OnCancel()```

* Get the result of future

  * ```.Get() ```
  * ```.GetOrTimeout()```
  * ```.GetChan()```

* Set timeout for future

  * ```.SetTimeout(ms)```
  
* Merge multiple promises

  * ```WhenAll(func1, func2, func3, ...)```
  * ```WhenAny(func1, func2, func3, ...)```
  * ```WhenAnyMatched(func1, func2, func3, ...)```

* Pipe
  * ```.Pipe(funcWithDone, funcWithFail)```

* Cancel the future

  * ```.Cancel()```
  * ```.IsCancelled()```

* Create future by function

  * ```Start(func() (r interface{}, e error))```
  * ```Start(func())```
  * ```Start(func(canceller Canceller) (r interface{}, e error))```
  * ```Start(func(canceller Canceller))```

* Immediate wrappers

  * ```Wrap(interface{})```

* Chain API

  * ```Start(taskDone).Done(done1).Fail(fail1).Always(alwaysForDone1).Pipe(f1, f2).Done(done2)```

## Quick start

### Promise and Future 

```go
import "github.com/fanliao/go-promise"
import "net/http"

p := promise.NewPromise()
p.OnSuccess(func(v interface{}) {
   ...
}).OnFailure(func(v interface{}) {
   ...
}).OnComplete(func(v interface{}) {
   ...
})

go func(){
	url := "http://example.com/"
	
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		p.Reject(err)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.Reject(err)
	}
	p.Resolve(body)
}()
r, err := p.Get()
```

If you want to provide a read-only view, you can get a future variable:

```go
p.Future //cannot Resolve, Reject for a future
```

Can use Start function to submit a future task, it will return a future variable, so cannot Resolve or Reject the future outside of Start function:

```go
import "github.com/fanliao/go-promise"
import "net/http"

task := func()(r interface{}, err error){
	url := "http://example.com/"
	
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

f := promise.Start(task).OnSuccess(func(v interface{}) {
   ...
}).OnFailure(func(v interface{}) {
   ...
}).OnComplete(func(v interface{}) {
   ...
})
r, err := f.Get()
```

### Get the result of future

Please note the process will be block until the future task is completed

```go
f := promise.Start(func() (r interface{}, err error) {
	return "ok", nil  
})
r, err := f.Get()  //return "ok", nil

f := promise.Start(func() (r interface{}, err error) {
	return nil, errors.New("fail")  
})
r, err := f.Get()  //return nil, errorString{"fail"}
```

Can wait until timeout

```go
f := promise.Start(func() (r interface{}, err error) {
	time.Sleep(500 * time.Millisecond)
	return "ok", nil 
})
r, err, timeout := f.GetOrTimeout(100)  //return nil, nil, true
```

### Merge multiple futures

Creates a future that will be completed when all of the supplied future are completed.
```go
task1 := func() (r interface{}, err error) {
	return "ok1", nil
}
task2 := func() (r interface{}, err error) {
	return "ok2", nil
}

f := promise.WhenAll(task1, task2)
r, err := f.Get()    //return []interface{}{"ok1", "ok2"}
```

If any future is failure, the future returnd by WhenAll will be failure
```go
task1 := func() (r interface{}, err error)  {
	return "ok", nil
}
task2 := func() (r interface{}, err error)  {
	return nil, errors.New("fail2")
}
f := promise.WhenAll(task1, task2)
r, ok := f.Get()    //return nil, *AggregateError
```

Creates a future that will be completed when any of the supplied tasks is completed.
```go
task1 := func() (r interface{}, err error) {
	return "ok1", nil
}
task2 := func() (r interface{}, err error) {
	time.Sleep(200 * time.Millisecond)
	return nil, errors.New("fail2")
}

f := promise.WhenAny(task1, task2)
r, err := f.Get()  //return "ok1", nil
```

Also can add a predicate function by WhenAnyMatched, the future that will be completed when any of the supplied tasks is completed and match the predicate.
```go
task1 := func() (r interface{}, err error) {
	time.Sleep(200 * time.Millisecond)
	return "ok1", nil
}
task2 := func() (r interface{}, err error) {
	return "ok2", nil
}

f := promise.WhenAnyMatched(func(v interface{}) bool{
	return v == "ok1"
}, task1, task2)
r, err := f.Get()  //return "ok1", nil
```

### Promise pipelining

```go
task1 := func() (r interface{}, err error) {
	return 10, nil
}
task2 := func(v interface{}) (r interface{}, err error) {
	return v.(int) * 2, nil
}

f := promise.Start(task1).Pipe(task2)
r, err := f.Get()   //return 20
```

### Cancel the future or set timeout

If need cancel a future, can pass a canceller object to task function
```go
import "github.com/fanliao/go-promise"
import "net/http"

p := promise.NewPromise().EnableCanceller()

go func(canceller promise.Canceller){
	for i < 50 {
		if canceller.IsCancelled() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}(p.Canceller())
f.Cancel()

r, err := p.Get()   //return nil, promise.CANCELLED
fmt.Println(p.Future.IsCancelled())      //true
```

Or can use Start to submit a future task which can be cancelled
```go
task := func(canceller promise.Canceller) (r interface{}, err error) {
	for i < 50 {
		if canceller.IsCancelled() {
			return 0, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return 1, nil
}
f := promise.Start(task1)
f.Cancel()

r, err := f.Get()   //return nil, promise.CANCELLED
fmt.Println(f.IsCancelled())      //true
```

When call WhenAny() function, if a future is completed correctly, then will try to check if other futures  enable cancel. If yes, will request cancelling all other futures.

You can also set timeout for a future

```go
task := func(canceller promise.Canceller) (r interface{}, err error) {
	time.Sleep(300 * time.Millisecond)
	if !canceller.IsCancelled(){
		fmt.Println("Run done")
	} 
	return
}

f := promise.Start(task).OnCancel(func() {
	fmt.Println("Future is cancelled")
}).SetTimeout(100)

r, err := f.Get() //return nil, promise.CANCELLED
fmt.Println(f.IsCancelled()) //print true
```

## Document

* [GoDoc at godoc.org](http://godoc.org/github.com/fanliao/go-promise)

## License

go-promise is licensed under the MIT Licence, (http://www.apache.org/licenses/LICENSE-2.0.html).
