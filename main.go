package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {

	// fmt.Println(runtime.NumCPU())

	// chanOwnerChanConsumer()
	// unsafeByteBufferConfinement()

	// parentSignalsSingleChildDone()
	// parentSignalsSingleChildDone2()
	errorHandling()

	// fiveCycles()

	// timeout()

}

func errorHandling() {

	type Result struct {
		Error    error
		Response *http.Response
	}
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

}

func parentSignalsSingleChildDone2() {

	// The previous example handles the case for goroutines receiving on a channel nicely,
	// but what if we’re dealing with the reverse situation: a goroutine blocked on attempting
	// to write a value to a channel? Here’s a quick example to demonstrate the issue:

	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}

func parentSignalsSingleChildDone() {

	//from Concurrency in Go
	//ch4

	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})

	strings := make(chan string)

	terminated := doWork(done, strings)

	go func() {
		// Cancel the operation after 1 second.

		//I added this bit. It is this or the sleep
		// for i := 0; i <= 9999; i++ {
		// 	strings <- fmt.Sprintf("%da", i)
		// }

		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

func chanOwnerChanConsumer() {
	//from Concurrency in Go
	//ch4 - confinement - 1
	//this format keeps the responsibilities of the 2 roles

	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner()
	consumer(results)

}

func select1() {

	//style 1 keeps the select short and does stuff after

	done := make(chan interface{})
	for {
		select {
		case <-done:
			return
		default:
		}

		// Do non-preemptable work
	}
}

func select2() {
	//style 2 embeds non-preemptable work in the default

	done := make(chan interface{})
	for {
		select {
		case <-done:
			return
		default:
			// Do non-preemptable work
		}
	}

}

func unsafeByteBufferConfinement() {

	//from Concurrency in Go
	//ch4 - confinement - 2

	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()

}

func fiveCycles() {
	done := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

func timeout() {

	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

}

func produceReceive() {

	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")

}

func main4() {

	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	// Increment
	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// Decrement
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")

}

func main2() {

	//1
	// var wg sync.WaitGroup
	// for _, salutation := range []string{"hello", "greetings", "good day"} {
	// 	wg.Add(1)
	// 	go func(salutation string) {
	// 		defer wg.Done()
	// 		fmt.Println(salutation)
	// 	}(salutation)
	// }
	// wg.Wait()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	wg.Wait()
	fmt.Println("All goroutines complete.")

}
