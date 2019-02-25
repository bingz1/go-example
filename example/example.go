package main

import (
	"fmt"
	"sync"
	"runtime"
)

var counter int = 0

func Count(lock *sync.Mutex) {
	lock.Lock() // 上锁
	counter++
	fmt.Println("counter =", counter)
	lock.Unlock()   // 解锁
}

func main() {
	lock := &sync.Mutex{}

	for i:=0; i<10; i++ {
		go Count(lock)
	}
	for {
		lock.Lock() // 上锁
		c := counter
		lock.Unlock()   // 解锁

		runtime.Gosched() // 出让时间片

		if c >= 10 {
			break
		}
	}
}