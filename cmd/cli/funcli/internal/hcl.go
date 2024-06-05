package internal

import (
	"fmt"
	"strings"
)

// MarshalArray marshals a slice of strings into a string that can be used in an HCL file.
func MarshalArray(s []string) string {
	cp := make([]string, len(s))
	copy(cp, s)

	for i := range cp {
		cp[i] = fmt.Sprintf("\"%v\"", cp[i])
	}

	return fmt.Sprintf("[%v]", strings.Join(cp, ", "))
}
