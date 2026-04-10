package main

import (
	"fmt"
	"io"
	"main/textprocessor"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("usage: go run . <input file> <output file>")
	fmt.Println("Options:")
	fmt.Println("  -h, --help	print this message and exit")
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printUsage()
		return
	}
	if len(os.Args) != 3 {
		printUsage()
		os.Exit(1)
	}

	inpath, outpath := os.Args[1], os.Args[2]

	var fileReader io.Reader
	if inpath == outpath {
		contents, err := os.ReadFile(inpath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "file \"%s\" not found", inpath)
				os.Exit(1)
			} else {
				fmt.Fprintln(os.Stderr, "Unexpected error:", err.Error())
				os.Exit(1)
			}
		}
		fileReader = strings.NewReader(string(contents))
	} else {
		infile, err := os.Open(inpath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "file \"%s\" not found", inpath)
				os.Exit(1)
			} else {
				fmt.Fprintln(os.Stderr, "Unexpected error:", err.Error())
				os.Exit(1)
			}
		}
		defer infile.Close()
		fileReader = infile
	}

	outfile, err := os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unexpected error:", err.Error())
		os.Exit(1)
	}
	defer outfile.Close()

	err = textprocessor.Run(fileReader, outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while processing:", err.Error())
		os.Remove(outpath)
		os.Exit(1)
	}
}
