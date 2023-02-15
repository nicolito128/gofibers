package main

import (
	"fmt"
	"sync"

	"github.com/nicolito128/gofibers"
)

func main() {
	// Start a WaitGroup in the main thread
	wg := &sync.WaitGroup{}
	// Create the new Fiber
	f := gofibers.New(wg, func(v string, r chan string) {
		fmt.Println("Initial message from the fiber:", v)
		gofibers.Suspend("Suspend 1", wg, r)

		fmt.Println("More code!")
	})
	defer f.Close()

	err := f.Start("Starting...")
	if err != nil {
		panic(err)
	}

	res, _ := f.Resume()
	// Printing the suspend message
	fmt.Println(res)

	fmt.Println("* Message outside the fiber.")

	// The last resume is an empty string
	f.Resume()
}
