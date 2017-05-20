package main

import (
	"os"
	"github.com/mgurov/hiptee/pkg"
	"github.com/mgurov/hiptee/pkg/std"
)

func main() {
	command := os.Args[1]
	params := os.Args[2:]
	if err := pkg.Execute(command, params, &std.StdOutPrinter{}, os.Stdin); nil != err {
		os.Exit(1)
	}

}