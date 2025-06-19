package engine

import (
	"fmt"
	"runtime"
	"sync"
)

type function[T any, R any] func(T) R

func parallelize[T any, R any](f function[T, R], args <-chan T) <-chan R {
	fmt.Println("Parallelize start")
	defer fmt.Println("Parallelize end")

	var numWorkers int = runtime.NumCPU()
	out := make(chan R, numWorkers*5)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for arg := range args {
				fmt.Printf("P : %v\n", arg)
				out <- f(arg)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
