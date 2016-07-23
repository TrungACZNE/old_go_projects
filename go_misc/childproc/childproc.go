package childproc

import (
	"bufio"
	"io"
	"log"
	"os/exec"

	"time"
)

func readerRoutine(reader *bufio.Reader, ch chan string, errCh chan error) {
loop:
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			break loop
		}
		ch <- s
	}
}

func watchFileWithTimeout(f io.ReadCloser, ch chan string, doneCh chan int, endChan chan int) {
	reader := bufio.NewReader(f)
	errCh := make(chan error)
	readCh := make(chan string)
	go readerRoutine(reader, readCh, errCh)
loop:
	for {
		select {
		case <-endChan:
			break loop
		case s := <-readCh:
			ch <- s
		case err := <-errCh:
			log.Println(err)
			break loop
		}
	}
	doneCh <- 1
	return
}

func watchFile(f io.ReadCloser, ch chan string, doneCh chan int) {
	reader := bufio.NewReader(f)
	errCh := make(chan error)
	readCh := make(chan string)
	go readerRoutine(reader, readCh, errCh)
loop:
	for {
		select {
		case s := <-readCh:
			ch <- s
		case err := <-errCh:
			log.Println(err)
			break loop
		}
	}
	doneCh <- 1
	return
}

// Start a child process which will write to childStdout and childStderr (line delimited)
// Any error produced during its life cycle will be written to errCh
// If timeoutSeconds > 0 the child process will receive a sigterm after that many seconds
func Start(commandName string, params []string, childStdout chan string, childStderr chan string, end chan int, errCh chan error, timeoutSeconds int) {
	cmd := exec.Command(commandName, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errCh <- err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		errCh <- err
	}

	err = cmd.Start()
	if err != nil {
		errCh <- err
	}

	doneChan := make(chan int)

	if timeoutSeconds > 0 {
		endChan := make(chan int)
		go watchFileWithTimeout(stdout, childStdout, doneChan, endChan)
		go watchFileWithTimeout(stderr, childStderr, doneChan, endChan)
		go func() {
			ticker := time.NewTicker(time.Duration(timeoutSeconds) * time.Second)
			<-ticker.C
			endChan <- 1
			endChan <- 1
			<-doneChan
			<-doneChan
			// Has to kill the process if it is still running
			// TODO verify that this won't blow up the planet
			err := cmd.Process.Kill()
			if err != nil {
				log.Println(err)
			}

		}()
	} else {
		go watchFile(stdout, childStdout, doneChan)
		go watchFile(stderr, childStderr, doneChan)
		<-doneChan
		<-doneChan
	}

	err = cmd.Wait()
	if err != nil {
		errCh <- err
	}
	end <- 1
}
