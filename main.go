package main

import (
	"fmt"
	"main/textprocessor"
	"os"
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

	outfile, err := os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unexpected error:", err.Error())
		os.Exit(1)
	}

	err = textprocessor.Run(infile, outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while processing:", err.Error())
		outfile.Close()
		os.Remove(outpath)
		os.Exit(1)
	}
}
