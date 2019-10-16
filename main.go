package main

import (
	"fmt"
	"log"
	"os"

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

	scriptRaw, readErr := read.SQLScript(outputPath)
	if readErr != nil {
		log.Panic(readErr)
	}
	fmt.Println(scriptRaw)
}
