package main

import (
	"github.com/mgurov/hiptee/pkg/hip"
	"io"
	"log"
	"os"
)

func main() {
	reader, err := hip.NewHipchatRoomReader(os.Getenv("HIPCHAT_TOKEN"), os.Getenv("HIPCHAT_ROOM"), "++")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, reader)
}
