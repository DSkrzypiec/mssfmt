package main

import (
	"fmt"
	"log"
	"os"

	"mssfmt/format"
	"mssfmt/read"
)

func main() {
	var outputPath string
	args := os.Args[1:]

	if len(args) < 1 {
		log.Panic("SQL script name or path is needed.")
	}
	outputPath = args[0]

	if len(args) >= 2 {
		outputPath = args[1]
	}

	ReadFormatAndSave(args[0], outputPath)
}

func ReadFormatAndSave(pathRead, pathWrite string) {
	scriptRaw, readErr := read.SQLScript(pathRead)
	if readErr != nil {
		log.Panic(readErr)
	}

	script := scriptRaw.ToScript()
	formats := format.BuildFormatsRepo()
	format.ApplyFormats(&script, formats)

	fmt.Println(script.RawContent)
	fmt.Println(script)
}
