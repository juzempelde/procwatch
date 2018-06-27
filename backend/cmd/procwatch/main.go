package main

import (
	"fmt"
	"os"
)

func main() {
	if err := createApp().Run(os.Args); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}

// Runner abstracts an application.
type Runner interface {
	Run([]string) error
}
