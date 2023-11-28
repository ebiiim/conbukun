package main

import (
	"fmt"

	"github.com/ebiiim/conbukun/pkg/handlers"
)

func main() {
	c := handlers.NewMapNameCompleter(8)

	// loop
	for {
		var input string
		fmt.Printf("input: ")
		if _, err := fmt.Scanln(&input); err != nil {
			continue
		}
		fmt.Printf("suggestions: %v\n", c.GetSuggestions(input))
	}

}
