package common

import (
	"fmt"
	"os"
)

func AssertErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Printf(msg+"\n", a...)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
