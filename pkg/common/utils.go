package common

import (
	"fmt"
	"os"
)

func AssertErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
