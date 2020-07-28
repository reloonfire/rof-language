package helpers

import "fmt"

// ReportError - Print error
func ReportError(line int, errMessage string) {
	fmt.Printf("[line %v] Error: %s\n", line, errMessage)
}
