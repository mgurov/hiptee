package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	count := flag.Int("count", 2, "how much ticks to do, one to stdout, another to sterr")
	exit := flag.Int("exit", 0, "exit code")
	delay := flag.Int("delay", 0, "delay, msec")
	flag.Parse()
	for i := 1; i <= *count; i++ {
		if i%2 == 0 {
			log.Println(i)
		} else {
			fmt.Println(i)
		}

		if *delay > 0 {
			time.Sleep(time.Duration(*delay) * time.Millisecond)
		}
	}
	os.Exit(*exit)
}
