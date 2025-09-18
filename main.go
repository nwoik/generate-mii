package main

import (
	"os"

	r "github.com/nwoik/Generate-Mii/rkg"
)

func main() {
	args := os.Args

	rkgFile := args[1]

	r.ExportToJsonRaw(rkgFile)
	r.ExportToJsonReadable(rkgFile)
	r.ExportMii(rkgFile)
}
