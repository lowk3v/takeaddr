package utils

import (
	"fmt"
	"os"
)

func HandleError(err error, customMsg string) bool {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", customMsg, err)
		return true
	}
	return false
}
