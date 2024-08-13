package main

import (
	"fmt"

	"github.com/nicolito128/gofibers"
)

func main() {
	// Create the new Fiber
	f := gofibers.New()
	// Use a fiber handler
	f.Handle(func(msg any) {
		fmt.Println("Initial message from the fiber:", msg)
		f.Suspend("Suspend 1")
		fmt.Println("More code!")
		f.Close()
	})
	// Send the initial message
	err := f.Init("Starting...")
	if err != nil {
		panic(err)
	}
	// Main goroutine work
	fmt.Println("Main goroutine!")
	// Receive the suspend
	res, err := f.Resume()
	if err != nil {
		panic(err)
	}
	fmt.Println(res.(string)) // Suspend 1
	// Wait for fiber close signal
	<-f.Closed()
}
