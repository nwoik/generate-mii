package main

import (
	"os"
	"strings"
	"time"

	r "github.com/nwoik/Generate-Mii/rkg"
)

func main() {
	args := os.Args

	rkgFile := args[1]

	if !strings.HasSuffix(rkgFile, ".rkg") {
		println("This isn't an rkg file")
		time.Sleep(3 * time.Second)
		return
	}

	r.ExportToJsonRaw(rkgFile)
	r.ExportToJsonReadable(rkgFile)
	r.ExportMii(rkgFile)
}
