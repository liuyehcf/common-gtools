package async

import (
	"github.com/liuyehcf/common-gtools/utils"
	"sync"
	"testing"
	"time"
)

func TestCond(t *testing.T) {
	m := &sync.Mutex{}
	c := sync.NewCond(m)
	m.Lock()
	go func() {
		logger.Info("try to get lock")
		m.Lock() // Wait for c.Wait()
		logger.Info("hold lock")
		c.Broadcast()
		m.Unlock()
	}()

	time.Sleep(time.Second)
	c.Wait() // Unlocks m, waits, then locks m again
	m.Unlock()
}

func TestRepeat(t *testing.T) {
	var promise Promise
	promise = CreatePromise()
	utils.AssertTrue(promise.TrySuccess(nil), "test")
	utils.AssertFalse(promise.TrySuccess(nil), "test")
	utils.AssertFalse(promise.TryFailure(nil), "test")
	utils.AssertFalse(promise.TryCancel(), "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")

	promise = CreatePromise()
	utils.AssertTrue(promise.TryFailure(nil), "test")
	utils.AssertFalse(promise.TrySuccess(nil), "test")
	utils.AssertFalse(promise.TryFailure(nil), "test")
	utils.AssertFalse(promise.TryCancel(), "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")

	promise = CreatePromise()
	utils.AssertTrue(promise.TryCancel(), "test")
	utils.AssertFalse(promise.TrySuccess(nil), "test")
	utils.AssertFalse(promise.TryFailure(nil), "test")
	utils.AssertFalse(promise.TryCancel(), "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertTrue(promise.IsCanceled(), "test")
}

func TestSync1(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TrySuccess(nil), "test")

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync2(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTrySuccess(promise, 300)
	}()

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync3(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTrySuccess(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			promise.Sync()

			defer lock.Unlock()
			lock.Lock()
			count++
		}()
	}

	for ; count < 10; {
	}
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync4(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TryFailure(nil), "test")

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync5(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryFailure(promise, 300)
	}()

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync6(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryFailure(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			promise.Sync()

			defer lock.Unlock()
			lock.Lock()
			count++
		}()
	}

	for ; count < 10; {
	}
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestSync7(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TryCancel(), "test")

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertTrue(promise.IsCanceled(), "test")
}

func TestSync8(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryCancel(promise, 300)
	}()

	promise.Sync()
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertTrue(promise.IsCanceled(), "test")
}

func TestSync9(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryCancel(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			promise.Sync()

			defer lock.Unlock()
			lock.Lock()
			count++
		}()
	}

	for ; count < 10; {
	}
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertFalse(promise.IsSuccess(), "test")
	utils.AssertTrue(promise.IsFailure(), "test")
	utils.AssertTrue(promise.IsCanceled(), "test")
}

func TestGet1(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TrySuccess("test"), "test")

	result := promise.Get()
	end := timestamp()

	utils.AssertNotNil(result, "test")
	utils.AssertTrue(end-start < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestGet2(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTrySuccess(promise, 300)
	}()

	result := promise.Get()
	end := timestamp()

	utils.AssertNotNil(result, "test")
	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestGet3(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTrySuccess(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			result := promise.Get()
			utils.AssertNotNil(result, "test")

			defer lock.Unlock()
			lock.Lock()
			count++
		}()
	}

	for ; count < 10; {
	}
	end := timestamp()

	utils.AssertTrue(end-start-300 < 100, "test")
	utils.AssertTrue(promise.IsDone(), "test")
	utils.AssertTrue(promise.IsSuccess(), "test")
	utils.AssertFalse(promise.IsFailure(), "test")
	utils.AssertFalse(promise.IsCanceled(), "test")
}

func TestGet4(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TryFailure(nil), "test")

	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")

		end := timestamp()
		utils.AssertTrue(end-start < 100, "test")
		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertFalse(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertFalse(promise.IsCanceled(), "test")
	}()
	promise.Get()
}

func TestGet5(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryFailure(promise, 300)
	}()

	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")

		end := timestamp()
		utils.AssertTrue(end-start-300 < 100, "test")
		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertFalse(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertFalse(promise.IsCanceled(), "test")
	}()
	promise.Get()
}

func TestGet6(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryFailure(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				err := recover()
				utils.AssertNotNil(err, "test")

				end := timestamp()
				utils.AssertTrue(end-start-300 < 100, "test")
				utils.AssertTrue(promise.IsDone(), "test")
				utils.AssertFalse(promise.IsSuccess(), "test")
				utils.AssertTrue(promise.IsFailure(), "test")
				utils.AssertFalse(promise.IsCanceled(), "test")

				defer lock.Unlock()
				lock.Lock()
				count++
			}()
			promise.Get()
		}()
	}

	for ; count < 10; {
	}
}

func TestGet7(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	utils.AssertTrue(promise.TryCancel(), "test")

	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")

		end := timestamp()
		utils.AssertTrue(end-start < 100, "test")
		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertFalse(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertTrue(promise.IsCanceled(), "test")
	}()
	promise.Get()
}

func TestGet8(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryCancel(promise, 300)
	}()

	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")

		end := timestamp()
		utils.AssertTrue(end-start-300 < 100, "test")
		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertFalse(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertTrue(promise.IsCanceled(), "test")
	}()
	promise.Get()
}

func TestGet9(t *testing.T) {
	promise := CreatePromise()
	start := timestamp()

	go func() {
		sleepThenTryCancel(promise, 300)
	}()

	lock := sync.Mutex{}
	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				err := recover()
				utils.AssertNotNil(err, "test")

				end := timestamp()
				utils.AssertTrue(end-start-300 < 100, "test")
				utils.AssertTrue(promise.IsDone(), "test")
				utils.AssertFalse(promise.IsSuccess(), "test")
				utils.AssertTrue(promise.IsFailure(), "test")
				utils.AssertTrue(promise.IsCanceled(), "test")

				defer lock.Unlock()
				lock.Lock()
				count++
			}()
			promise.Get()
		}()
	}

	for ; count < 10; {
	}
}

func TestListener1(t *testing.T) {
	var isTriggered bool
	var promise Promise
	isTriggered = false
	promise = CreatePromise()
	utils.AssertTrue(promise.TrySuccess(nil), "test")
	promise.AddListener(func(promise Promise) {
		isTriggered = true

		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertTrue(promise.IsSuccess(), "test")
		utils.AssertFalse(promise.IsFailure(), "test")
		utils.AssertFalse(promise.IsCanceled(), "test")
	})
	utils.AssertTrue(isTriggered, "test")

	isTriggered = false
	promise = CreatePromise()
	utils.AssertTrue(promise.TryFailure(nil), "test")
	promise.AddListener(func(promise Promise) {
		isTriggered = true

		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertTrue(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertFalse(promise.IsCanceled(), "test")
	})
	utils.AssertTrue(isTriggered, "test")

	isTriggered = false
	promise = CreatePromise()
	utils.AssertTrue(promise.TryCancel(), "test")
	promise.AddListener(func(promise Promise) {
		isTriggered = true

		utils.AssertTrue(promise.IsDone(), "test")
		utils.AssertFalse(promise.IsSuccess(), "test")
		utils.AssertTrue(promise.IsFailure(), "test")
		utils.AssertTrue(promise.IsCanceled(), "test")
	})
	utils.AssertTrue(isTriggered, "test")
}

func TestListener2(t *testing.T) {
	promise := CreatePromise()
	count1 := 0
	promise.AddListener(func(promise Promise) {
		count1++
	})

	utils.AssertTrue(promise.TrySuccess(nil), "nil")

	count2 := 0
	promise.AddListener(func(promise Promise) {
		count2++
	})

	utils.AssertTrue(count1 == 1, "test")
	utils.AssertTrue(count2 == 1, "test")
}

func TestListener3(t *testing.T) {
	promise := CreatePromise()

	count := 0
	lock := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		promise.AddListener(func(promise Promise) {
			defer lock.Unlock()
			lock.Lock()
			count++
			utils.AssertTrue(promise.IsDone(), "test")
			utils.AssertTrue(promise.IsSuccess(), "test")
			utils.AssertFalse(promise.IsFailure(), "test")
			utils.AssertFalse(promise.IsCanceled(), "test")
		})
	}

	utils.AssertTrue(promise.TrySuccess(nil), "nil")

	utils.AssertTrue(count == 1000, "test")
}

func TestFailed(t *testing.T) {
	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")
		errContent := err.(string)
		utils.AssertTrue("promise failed" == errContent, "test")
	}()

	promise := CreatePromise()
	promise.AddListener(func(promise Promise) {
		cause := promise.Cause()
		causeContent := cause.(string)
		utils.AssertTrue("test failure" == causeContent, "test")
	})

	promise.TryFailure("test failure")

	promise.Get()
}

func TestCancel(t *testing.T) {
	defer func() {
		err := recover()
		utils.AssertNotNil(err, "test")
		s := err.(string)
		utils.AssertTrue("promise canceled" == s, "test")
	}()

	promise := CreatePromise()

	promise.TryCancel()

	promise.Get()
}

func timestamp() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func sleepThenTrySuccess(promise Promise, milliseconds time.Duration) {
	time.Sleep(time.Millisecond * milliseconds)
	utils.AssertTrue(promise.TrySuccess("test"), "test")
}

func sleepThenTryFailure(promise Promise, milliseconds time.Duration) {
	time.Sleep(time.Millisecond * milliseconds)
	utils.AssertTrue(promise.TryFailure("test"), "test")
}

func sleepThenTryCancel(promise Promise, milliseconds time.Duration) {
	time.Sleep(time.Millisecond * milliseconds)
	utils.AssertTrue(promise.TryCancel(), "test")
}
