package proxy

import "os"

func AutoSet() {
	os.Setenv("GO111MODULE", "on")
	os.Setenv("GOPROXY", "https://goproxy.io")
}
