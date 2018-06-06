package main

import "fmt"

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}

func bridge(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	// Channel that returns all values from bridge
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		// Loop responsible for pulling channels off chanStream
		// and providing them to a nested loop for use.
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			// Responsible for reading values off the channel given
			// and repeating those onto valStream.
			// When stream we're currently looping over is closed, we break out
			// of the loop performing the reads form this channel, and continue
			// with next iteration fo the loop, selecting channels to read from.
			// This provides us with an unbroken stream of values.
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
