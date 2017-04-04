package std

import (
	"fmt"
	"log"
	"os"
)

type StdOutPrinter struct {
}

func (s *StdOutPrinter) Out(line string) {
	fmt.Fprintln(os.Stdout, line)
}

func (s *StdOutPrinter) Err(line string) {
	fmt.Fprintln(os.Stderr, line)
}

func (s *StdOutPrinter) Done(err error) {
	if err != nil {
		log.Println("Error executing subprocess", err)
	}
}
