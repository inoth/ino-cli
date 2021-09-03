package proxy

import "os"

func AutoSet() {
	SetGoModuleEnabled(true)
	os.Setenv("GOPROXY", "https://goproxy.io")
}

func SetGoModuleEnabled(enabled bool) {
	if enabled {
		os.Setenv("GO111MODULE", "on")
	} else {
		os.Setenv("GO111MODULE", "off")
	}
}
