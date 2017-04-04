package main

import (
	"github.com/mgurov/teehip/pkg"
	"github.com/mgurov/teehip/pkg/std"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Need an argument to run ", os.Args)
	}
	command := os.Args[1]
	params := os.Args[2:]
	log.Println("About to execute", command, "with params", params)

	if err := pkg.Execute(command, params, new(std.StdOutPrinter)); nil != err {
		os.Exit(1)
	}
}
