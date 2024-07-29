package fins

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_singleflightOne_do(t *testing.T) {
	var sg singleflightOne
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			fmt.Println(recover())
			wg.Done()
		}()
		cc(&sg)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			fmt.Println(recover())
			wg.Done()
		}()
		cc(&sg)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			fmt.Println(recover())
			wg.Done()
		}()
		cc(&sg)
	}()

	wg.Wait()
}

func cc(sg *singleflightOne) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sg.do(f2)
	}()
	sg.do(f1)
}

func f1() {
	time.Sleep(time.Second)
	panic(1)
}

func f2() {
	time.Sleep(time.Second)
	fmt.Println(1)
}
