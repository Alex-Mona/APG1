package main

import (
	"fmt"
	"sync"
)

// multiplex принимает несколько каналов и объединяет их в один
func multiplex(channels ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c chan interface{}) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go func() {
		ch1 <- "from ch1"
		close(ch1)
	}()

	go func() {
		ch2 <- "from ch2"
		close(ch2)
	}()

	out := multiplex(ch1, ch2)

	// Читаем из объединённого канала
	for val := range out {
		fmt.Println(val)
	}
}
