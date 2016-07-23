package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

type Blob interface {
	Value() string
}

type RawBlob string

func (blob RawBlob) Value() string {
	return string(blob)
}

type PermBlob []string

func (blob PermBlob) Value() string {
	org := []string(blob)
	return org[rand.Int()%len(org)]
}

func inputError(reason, text string, pos int) error {
	row, col := getPos(text, pos)
	return fmt.Errorf("%s at %d:%d\n", reason, row, col)
}

func getPos(text string, i int) (int, int) {
	if i >= len(text) {
		i = len(text)
	}
	t := text[:i]
	newLines := strings.Count(t, "\n")
	lastNewLine := strings.LastIndex(t, "\n")
	return newLines + 1, i - lastNewLine
}

func parsePermutation(text string, start int) ([]string, int, error) {
	l := len(text)
	choices := []string{}
	acc := ""
	i := start
	for ; i < l; i++ {
		if text[i] == '}' && (text[i-1] != '\\') {
			choices = append(choices, acc)
			return choices, i, nil
		} else if text[i] == ',' {
			choices = append(choices, acc)
			acc = ""
		} else {
			acc += string(text[i])
		}
	}
	return choices, i, inputError("Unexpected end of input", text, start)
}

func parse(text string) []Blob {
	result := []Blob{}
	i := 0
	l := len(text)
	acc := ""
	for ; i < l; i++ {
		if text[i] == '{' && (i == 0 || text[i-1] != '\\') {
			result = append(result, RawBlob(acc))
			acc = ""
			choices, last, err := parsePermutation(text, i+1)
			if err != nil {
				log.Fatal(err)
			}
			result = append(result, PermBlob(choices))
			i = last
		} else {
			acc += string(text[i])
		}
	}
	result = append(result, RawBlob(acc))
	return result
}

func generate(text string) string {
	blobs := parse(text)
	result := ""
	for _, blob := range blobs {
		result += blob.Value()
	}
	return result
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "Word shuffler"
	app.Usage = "Word shuffling tool. Reads from input writes to output. See \"email\" for an example."
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) {
		rand.Seed(time.Now().Unix())
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(generate(string(input)))
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("app.Run() error:", err)
	}
}
