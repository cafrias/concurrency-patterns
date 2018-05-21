package main

// Resource: Concurrency in Go, Ch. 4. Katherine Cox-Buday.

func main() {
	done := make(chan interface{})
	stringStream := make(chan string)

	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
}
