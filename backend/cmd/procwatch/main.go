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
