package async

import (
	"github.com/liuyehcf/common-gtools/log"
	"sync"
	"sync/atomic"
)

const (
	notTriggered = 0
	triggered    = 1
)

var logger = log.GetLogger("Promise")

type Promise interface {

	// return true if this task was canceled before it completed normally
	IsCanceled() bool

	// return true if this task completed
	IsDone() bool

	// return true if this task succeeded
	IsSuccess() bool

	// return true if this task failed
	IsFailure() bool

	// return cause of the failure
	Cause() interface{}

	// attempts to cancel execution of this task
	TryCancel() bool

	// marks this promise as a success and notifies all listeners if this promise isn't done
	TrySuccess(outcome interface{}) bool

	// marks this promise as a failure and notifies all listeners if this promise isn't done
	TryFailure(cause interface{}) bool

	// add the specified listener to this promise
	AddListener(listener func(promise Promise)) Promise

	// waits for this promise until it is done, and panic with the cause of the failure if this promise failed
	Sync()

	// waits if necessary for the computation to complete, and then retrieves its result
	Get() interface{}
}

type defaultPromise struct {
	listeners  []*promiseListenerWrapper
	lock       *sync.Mutex
	condition  *sync.Cond
	outcome    interface{}
	cause      interface{}
	isCanceled bool
	isDone     bool
	isSuccess  bool
	isFailure  bool
}

func CreatePromise() Promise {
	mutex := &sync.Mutex{}
	return &defaultPromise{
		listeners:  make([]*promiseListenerWrapper, 0),
		lock:       mutex,
		condition:  sync.NewCond(mutex),
		outcome:    nil,
		cause:      nil,
		isCanceled: false,
		isDone:     false,
		isSuccess:  false,
		isFailure:  false,
	}
}

func (promise *defaultPromise) IsCanceled() bool {
	return promise.isCanceled
}

func (promise *defaultPromise) IsDone() bool {
	return promise.isDone
}

func (promise *defaultPromise) IsSuccess() bool {
	return promise.isSuccess
}

func (promise *defaultPromise) IsFailure() bool {
	return promise.isFailure
}

func (promise *defaultPromise) Cause() interface{} {
	return promise.cause
}

func (promise *defaultPromise) TryCancel() bool {
	if promise.IsDone() {
		return false
	}

	// guarantee only one of three methods(setCanceledUnderLock, setSuccessUnderLock, setFailureUnderLock) can be execute
	// and only execute only once
	result := promise.executeSynchronousUnderLock(func() interface{} {
		if promise.IsDone() {
			return false
		}

		promise.setCanceledUnderLock()

		return true
	})

	promise.notifyAllListeners()

	return result.(bool)
}

func (promise *defaultPromise) TrySuccess(outcome interface{}) bool {
	if promise.IsDone() {
		return false
	}

	// guarantee only one of three methods(setCanceledUnderLock, setSuccessUnderLock, setFailureUnderLock) can be execute
	// and only execute only once
	result := promise.executeSynchronousUnderLock(func() interface{} {
		if promise.IsDone() {
			return false
		}

		promise.setSuccessUnderLock(outcome)

		return true
	})

	promise.notifyAllListeners()

	return result.(bool)
}

func (promise *defaultPromise) TryFailure(cause interface{}) bool {
	if promise.IsDone() {
		return false
	}

	// guarantee only one of three methods(setCanceledUnderLock, setSuccessUnderLock, setFailureUnderLock) can be execute
	// and only execute only once
	result := promise.executeSynchronousUnderLock(func() interface{} {
		if promise.IsDone() {
			return false
		}

		promise.setFailureUnderLock(cause)

		return true
	})

	promise.notifyAllListeners()

	return result.(bool)
}

func (promise *defaultPromise) AddListener(listener func(promise Promise)) Promise {
	promise.addListener0(listener)

	if promise.IsDone() {
		promise.notifyAllListeners()
	}

	return promise
}

func (promise *defaultPromise) Sync() {
	if promise.IsDone() {
		return
	}

	// guarantee when this current thread is going to block itself,
	// there must be other thread to wait it up
	promise.executeSynchronousUnderLock(func() interface{} {
		if promise.IsDone() {
			return nil
		}

		promise.condition.Wait()
		return nil
	})
}

func (promise *defaultPromise) Get() interface{} {
	if promise.IsDone() {
		return promise.report()
	}

	// guarantee when this current thread is going to block itself,
	// there must be other thread to wait it up
	return promise.executeSynchronousUnderLock(func() interface{} {
		if promise.IsDone() {
			return promise.report()
		}

		promise.condition.Wait()
		return promise.report()
	})
}

func (promise *defaultPromise) executeSynchronousUnderLock(callable func() interface{}) interface{} {
	defer promise.lock.Unlock()
	promise.lock.Lock()

	return callable()
}

func (promise *defaultPromise) setCanceledUnderLock() {
	// isDone is used to determine success, so the assignment of isDone must be at the end
	promise.isFailure = true
	promise.isCanceled = true
	promise.isDone = true
	promise.condition.Broadcast()
}

func (promise *defaultPromise) setSuccessUnderLock(outcome interface{}) {
	// isDone is used to determine success, so the assignment of isDone must be at the end
	promise.isSuccess = true
	promise.outcome = outcome
	promise.isDone = true
	promise.condition.Broadcast()
}

func (promise *defaultPromise) setFailureUnderLock(cause interface{}) {
	// isDone is used to determine success, so the assignment of isDone must be at the end
	promise.isFailure = true
	promise.cause = cause
	promise.isDone = true
	promise.condition.Broadcast()
}

func (promise *defaultPromise) report() interface{} {
	if promise.IsSuccess() {
		return promise.outcome
	}

	if promise.IsCanceled() {
		panic("promise canceled")
	}

	panic("promise failed")
}

func (promise *defaultPromise) addListener0(listener func(promise Promise)) {
	promise.executeSynchronousUnderLock(func() interface{} {
		promise.listeners = append(promise.listeners, &promiseListenerWrapper{
			target:      listener,
			isTriggered: notTriggered,
		})
		return nil
	})
}

func (promise *defaultPromise) notifyAllListeners() {
	for _, listener := range promise.listeners {
		promise.notifyListener(listener)
	}
}

func (promise *defaultPromise) notifyListener(listener *promiseListenerWrapper) {
	if atomic.CompareAndSwapInt32(&listener.isTriggered, notTriggered, triggered) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn("an exception was thrown by PromiseListener, errorMsg={}", err)
			}
		}()
		listener.target(promise)
	}
}

type promiseListenerWrapper struct {
	target      func(promise Promise)
	isTriggered int32
}
