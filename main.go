package main

import (
	"log"
	"os"

	"github.com/inoth/ino-cli/cmd"
	"github.com/inoth/ino-cli/proxy"
)

func init() {
	proxy.AutoSet()
}

func main() {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				log.Fatalf(err.Error())
			} else {
				panic(exception)
			}
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
