package util

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}
