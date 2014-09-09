package main

import (
	"fmt"
	"os"
	"./sql"
)

func showHelp() {
	fmt.Fprintf(os.Stderr, helpText)
}

func main() {
	if len(os.Args) < 3 {
		showHelp()
		os.Exit(0)
	}

	src, err := sql.GetSchema(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	dst, err := sql.GetSchema(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	err = src.GetSchema()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = dst.GetSchema()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	diff := sql.DiffSchema(src, dst)
	fmt.Println(diff)
}

const helpText = `Usage: mysqldiff <src db / sql path> <dist db / sql path>
`
