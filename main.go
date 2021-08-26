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
	// errorHandling()
	// pipelinesBatchStyle()
	// pipelinesStreamStyle()
	// pipelinesWithChannels()
	// repeatAndTakeGenerators()
	repeatFnGenerator()
	// fanInFanOutPrimes()

	// fiveCycles()

	// timeout()

}

func fanInFanOutPrimes() {

	//this is missing impls for primeFinder and toInt

	// take := func(
	// 	done <-chan interface{},
	// 	valueStream <-chan interface{},
	// 	num int,
	// ) <-chan interface{} {
	// 	takeStream := make(chan interface{})
	// 	go func() {
	// 		defer close(takeStream)
	// 		for i := 0; i < num; i++ {
	// 			select {
	// 			case <-done:
	// 				return
	// 			case takeStream <- <-valueStream:
	// 			}
	// 		}
	// 	}()
	// 	return takeStream
	// }

	// repeatFn := func(
	// 	done <-chan interface{},
	// 	fn func() interface{},
	// ) <-chan interface{} {
	// 	valueStream := make(chan interface{})
	// 	go func() {
	// 		defer close(valueStream)
	// 		for {
	// 			select {
	// 			case <-done:
	// 				return
	// 			case valueStream <- fn():
	// 			}
	// 		}
	// 	}()
	// 	return valueStream
	// }

	// fanIn := func(
	// 	done <-chan interface{},
	// 	channels ...<-chan interface{},
	// ) <-chan interface{} {
	// 	var wg sync.WaitGroup
	// 	multiplexedStream := make(chan interface{})

	// 	multiplex := func(c <-chan interface{}) {
	// 		defer wg.Done()
	// 		for i := range c {
	// 			select {
	// 			case <-done:
	// 				return
	// 			case multiplexedStream <- i:
	// 			}
	// 		}
	// 	}

	// 	// Select from all the channels
	// 	wg.Add(len(channels))
	// 	for _, c := range channels {
	// 		go multiplex(c)
	// 	}

	// 	// Wait for all the reads to complete
	// 	go func() {
	// 		wg.Wait()
	// 		close(multiplexedStream)
	// 	}()

	// 	return multiplexedStream
	// }

	// done := make(chan interface{})
	// defer close(done)

	// start := time.Now()

	// rand := func() interface{} { return rand.Intn(50000000) }

	// randIntStream := toInt(done, repeatFn(done, rand))

	// numFinders := runtime.NumCPU()
	// fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	// finders := make([]<-chan interface{}, numFinders)
	// fmt.Println("Primes:")
	// for i := 0; i < numFinders; i++ {
	// 	finders[i] = primeFinder(done, randIntStream)
	// }

	// for prime := range take(done, fanIn(done, finders...), 10) {
	// 	fmt.Printf("\t%d\n", prime)
	// }

	// fmt.Printf("Search took: %v", time.Since(start))

}

func repeatFnGenerator() {

	//an infinite channel of random integers generated on an as-needed basis!

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}

func repeatAndTakeGenerators() {

	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

}

func pipelinesWithChannels() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func pipelinesStreamStyle() {

	multiply := func(value, multiplier int) int {
		return value * multiplier
	}

	add := func(value, additive int) int {
		return value + additive
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
		fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	}

}

func pipelinesBatchStyle() {

	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	ints = []int{1, 2, 3, 4}
	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}

}

func errorHandling() {

	//this actually makes concurrent http requests?
	//can I achieve 10/100 with this code?

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

	errCount := 0
	urls := []string{"a", "https://www.google.com", "https://badhost", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
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
