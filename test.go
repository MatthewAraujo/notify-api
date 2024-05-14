package main

import (
	"fmt"
	"os"
)

func main() {
	// Define payload
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current working directory:", wd)
}
