package main

import (
	"fmt"
	"time"

	"sync"
)

func or(chs ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	stop := make(chan struct{})
	var once sync.Once

	for _, ch := range chs {
		go func(c <-chan interface{}) {
			select {
			case _, ok := <-c:
				if !ok {
					once.Do(func() {
						close(stop)
						close(done)
					})
				}
			case <-stop:
				return
			}
		}(ch)
	}

	return done
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(3*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
