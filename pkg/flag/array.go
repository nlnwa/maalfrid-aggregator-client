package flag

import "strings"

// An ArrayFlag holds the values for multiple flags with the same name
type ArrayFlag []string

// String returns the string representation of the flags in the array
func (arrayFlag *ArrayFlag) String() string {
	return strings.Join(*arrayFlag, ",")
}

// Set the value of the flag
func (arrayFlag *ArrayFlag) Set(value string) error {
	*arrayFlag = append(*arrayFlag, value)
	return nil
}
