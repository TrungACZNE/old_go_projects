package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	RE_INPUT := regexp.MustCompile(`(#[^\s]+)\s([^\s]+)\s(.*)`)
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: '%v'", err.Error())
			os.Exit(1)
		}

		if matches := RE_INPUT.FindStringSubmatch(line); len(matches) == 4 {
			channel := matches[1]
			username := matches[2]
			// message := matches[3]
			fmt.Printf("PRIVMSG %s :PogChamp %s\n", channel, username)
		} else {
			fmt.Println()
		}
	}
}
