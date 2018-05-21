// Lexical confinement
package main

// Resource: Concurrency in Go, Ch. 4. Katherine Cox-Buday.

import "fmt"

func main() {
	// By returning a read-only channel we ensure this channel cannot be written outside the scope of this function
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

	// consumer receives a read-only channel, just in case
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner()
	consumer(results)
}
