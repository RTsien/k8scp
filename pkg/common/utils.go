package common

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

func AssertErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type mutexReader struct {
	reader io.Reader
	mu     sync.Mutex
	caller string
}

func MutexReader(reader io.Reader) io.Reader {
	return &mutexReader{reader: reader}
}

func (r *mutexReader) Read(p []byte) (int, error) {
	if !r.mu.TryLock() {
		fmt.Fprintf(os.Stderr, "----------------reader conflict----------------\ncur: %s\nother: %s", r.callerInfo(), r.caller)
		r.mu.Lock()
	}
	defer r.mu.Unlock()

	r.caller = r.callerInfo()

	return r.reader.Read(p)
}

func (r *mutexReader) callerInfo() (info string) {
	for i := 2; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		pcName := runtime.FuncForPC(pc).Name()
		info += fmt.Sprintf("%v %s %d %t %s %s\n", pc, file, line, ok, pcName, time.Now())
	}
	return
}
