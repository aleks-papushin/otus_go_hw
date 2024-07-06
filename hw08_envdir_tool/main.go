package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "error: you should pass at least 2 arguments:\n"+
			"1) path to directory containing env variable files\n"+
			"2) path to program which should process files\n")
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Sprintf("ReadDir method failed with: %s", err.Error()))
		os.Exit(1)
	}

	RunCmd(os.Args[2:], env)
}
