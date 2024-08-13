# Go Fibers
Go Fibers is an attempt to implement <a href="https://www.php.net/manual/es/language.fibers.php">PHP Fibers</a> in the language. There are probably better ways to do this type of implementation.

## Requirements
* Go version `1.22+`

## Installation

Get the module

    go get -u github.com/nicolito128/gofibers

Then import it

	import "github.com/nicolito128/gofibers"

## Usage example

```go
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
```
