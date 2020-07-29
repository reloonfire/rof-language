package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		return
	}
	outputDir := os.Args[1]

}
