package flag

import "strings"

// An ArrayFlag holds the values for same named flags
type ArrayFlag []string

// String returns the string representation of the flags in the array
func (i *ArrayFlag) String() string {
	return strings.Join(*i, ",")
}

// Set the value of the flag
func (i *ArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
