package utils

import (
	"log"
	"runtime"
	"sync"

	uuid "github.com/odeke-em/go-uuid"
)

type LogLock struct {
	mu sync.RWMutex
}

func MyCaller() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	return fun.Name()
}

func (l *LogLock) Lock() {
	u := uuid.New()
	log.Println(u, MyCaller(), "lock wait")
	l.mu.Lock()
	log.Println(u, MyCaller(), "lock get")
}

func (l *LogLock) Unlock() {
	l.mu.Unlock()
}

func (l *LogLock) RLock() {
	u := uuid.New()
	log.Println(u, MyCaller(), "rlock wait")
	l.mu.RLock()
	log.Println(u, MyCaller(), "rlock get")
}

func (l *LogLock) RUnlock() {
	l.mu.RUnlock()
}
