package main

import (
	"ino-cli/cmd"
	"ino-cli/command/init"
	"ino-cli/proxy"
	"log"
	"strings"
)

const (
	varsion = "v1.0"
)

var (
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
	helpContent = strings.TrimLeft(`
USAGE
    ino-cli COMMAND [ARGUMENT] [OPTION]
COMMAND
	version    print version
    init       create and initialize an empty project...
`, DefaultTrimChars)
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
		}
	}()
	command := cmd.GetArg(1)
	switch command {
	case "help":
		help(cmd.GetArg(2))
	case "init":
		init.Run()
	case "version":
		version()
	default:
		log.Fatalln(helpContent)
	}
}

func help(command string) {
	switch command {
	case "init":
		init.Help()
	default:
		log.Fatalln(helpContent)
	}
}

func version() {
	log.Fatalln(varsion)
}
