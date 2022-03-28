package errhelp

import (
	"fmt"
	"os"
)

func ThrowConnectionError(err error) {
	fmt.Printf("Unexpected connection error: %v", err)
	os.Exit(3)
}

func ThrowCloseError(err error) {
	fmt.Printf("Unexpected error while closing: %v", err)
	os.Exit(3)
}
