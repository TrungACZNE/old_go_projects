package childproc

import (
	"log"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	stdout := make(chan string)
	stderr := make(chan string)
	end := make(chan int)
	err := make(chan error)
	go Start("python", []string{"test.py"}, stdout, stderr, end, err, 4)
loop:
	for {
		select {
		case s := <-stdout:
			log.Println("Received from stdout", strings.Trim(s, "\n"))
		case s := <-stderr:
			log.Println("Received from stderr", strings.Trim(s, "\n"))
		case <-end:
			log.Println("Finished")
			break loop
		case myerr := <-err:
			log.Println("Fucked up", myerr)
			break loop

		}
	}
}
