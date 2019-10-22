package main

import (
	"fmt"
	"log"
	"os"

	"mssfmt/read"
	"mssfmt/scanner"
	"mssfmt/token"
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

	var scriptScanner scanner.Scanner
	scriptScanner.Init(outputPath, []byte(scriptRaw.Content))

	for {
		tok, lit := scriptScanner.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("[%s] %s \n", tok, lit)
	}
}
