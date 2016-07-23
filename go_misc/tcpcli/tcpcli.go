package tcpcli

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/TrungACZNE/go_misc/cliparser"
)

type CommandHandler func(string, []string, error)

func Start(bind string, cmdHandler CommandHandler) error {
	l, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	defer l.Close()
	fmt.Println("Listening on " + bind)
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handleRequest(conn, cmdHandler)
	}
	return nil
}

// Handles incoming requests.
func handleRequest(conn net.Conn, cmdHandler CommandHandler) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmdHandler(cliparser.Parse(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
