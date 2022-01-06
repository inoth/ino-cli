package cmd

import (
	"os"
	"regexp"
)

var (
	argumentRegex     = regexp.MustCompile(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`)
	defaultParsedArgs = make([]string, 0)
)

func Init() {
	defaultParsedArgs = make([]string, 0)
	args := os.Args
	for i := 0; i < len(args); {
		array := argumentRegex.FindStringSubmatch(args[i])
		if len(array) > 2 {
		} else {
			defaultParsedArgs = append(defaultParsedArgs, args[i])
		}
		i++
	}
}

func GetArg(i int) string {
	Init()
	if i < len(defaultParsedArgs) {
		return defaultParsedArgs[i]
	}
	return ""
}
