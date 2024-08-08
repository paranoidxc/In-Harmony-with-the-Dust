package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/single_flight"
	"sync"
	"time"
)

var mu sync.Mutex
var num = 0

func get() (interface{}, error) {
	mu.Lock()
	num++
	e := num
	time.Sleep(3 * time.Second)
	mu.Unlock()
	return e, nil
}

func main() {
	wg := sync.WaitGroup{}
	sf := singleFlight.Group{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(sf.Do("echo", get))
			wg.Done()
		}()
		if i == 5 {
			time.Sleep(5 * time.Second)
			sf.ForgetUnshared("echo")
		}
	}
	wg.Wait()
	ch1 := sf.DoChan("echo", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 1, nil
	})
	ch2 := sf.DoChan("echo", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 2, nil
	})
	ch3 := sf.DoChan("echo", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 3, nil
	})
	select {
	case p := <-ch1:
		fmt.Println(p)
	case p := <-ch2:
		fmt.Println(p)
	case p := <-ch3:
		fmt.Println(p)
	case <-time.After(3 * time.Second):
		fmt.Println("超时")
	}
}
