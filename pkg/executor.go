package pkg

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

func Execute(command string, params []string, listener LineOutputListener, stdin io.Reader) error {
	cmd := exec.Command(command, params...)
	wg := sync.WaitGroup{}

	cmd.Stdin = stdin

	if err := listen(listener.Out, cmd.StdoutPipe, &wg); err != nil {
		log.Fatal(err)
	}

	if err := listen(listener.Err, cmd.StderrPipe, &wg); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	err := cmd.Wait()
	listener.Done(err)
	return err
}

func listen(listenerFun func(string), pipeProvider func() (io.ReadCloser, error), wg *sync.WaitGroup) error {
	pipe, err := pipeProvider()
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			listenerFun(scanner.Text())
		}
		wg.Done()
	}()
	return nil
}
