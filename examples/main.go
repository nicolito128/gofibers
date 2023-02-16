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
	f := gofibers.New(wg, func(v any, r chan any) {
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
	fmt.Println(res.(string)) // Suspend 1

	// Last handler execution
	f.Resume()
}
