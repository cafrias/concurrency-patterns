package main

// Resource: Concurrency in Go, Ch. 4. Katherine Cox-Buday.

func main() {
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
