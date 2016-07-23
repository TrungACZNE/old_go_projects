package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/TrungACZNE/go_twitch_bot/child_process"
)

func fail(location string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s -> Fatal error: %s", location, err.Error())
		os.Exit(1)
	}
}

const (
	TYPE_NIL  = iota
	TYPE_CHAT = iota
)

type Command struct {
	Username    string
	CommandName string
	Params      string

	Type        int
	ChatMessage string
	ChatChannel string
}

var (
	RE_COMMAND     *regexp.Regexp
	RE_CHAT_PARAMS *regexp.Regexp

	RE_SEARCH *regexp.Regexp

	CHAT_COMMAND = "PRIVMSG"

	username string
	password string
	channel  string

	ignoreAdmin bool

	searchReg        string
	replyCommand     string
	replyCommandLoop string
)

func init() {
	RE_COMMAND = regexp.MustCompile(`:([^!]+)[^\s]+\s([A-Z]+)\s(.*)`)
	RE_CHAT_PARAMS = regexp.MustCompile(`([^\s]+)\s:(.*)`)

	flag.StringVar(&username, "username", "", "Your twitch username")
	flag.StringVar(&password, "password", "", "Your oauth password generated at http://twitchapps.com/tmi/")
	flag.StringVar(&channel, "channel", "", "The twitch channel to monitor")

	flag.BoolVar(&ignoreAdmin, "ignoreAdmin", false, "Ignores message from jtv (which tells whether your message was sent successfully or not)")
	flag.StringVar(&searchReg, "search", "", "Only prints chat messages that match this Golang regexp string")
	flag.StringVar(&replyCommand, "exec", "", "Used with --search, when a match is found execute the program named by this parameter with parameters (channel name, username, password). Channel name IS prefixed with the # sign. The output of this program will be sent as an IRC command back to the channel.")
	flag.StringVar(&replyCommandLoop, "execLoop", "", "Similar to --exec but the command is supposed to work like a shell")
}

func (command *Command) ResolveChatMessage() bool {
	matches := RE_CHAT_PARAMS.FindStringSubmatch(command.Params)
	if len(matches) == 3 {
		command.Type = TYPE_CHAT
		command.ChatChannel = matches[1]
		command.ChatMessage = matches[2]
		return true
	}
	return false
}

func StrToCommand(line string) *Command {
	matches := RE_COMMAND.FindStringSubmatch(line)
	if len(matches) == 4 {
		v := &Command{}
		v.Username = matches[1]
		v.CommandName = matches[2]
		v.Params = matches[3]

		v.Type = TYPE_NIL
		v.ChatChannel = ""
		v.ChatMessage = ""

		v.ResolveChatMessage()
		return v
	} else {
		return nil
	}
}

func main() {
	flag.Parse()

	if username == "" || password == "" || channel == "" {
		fmt.Fprintf(os.Stderr, "Missing either --username, --password or --channel")
		os.Exit(1)
	}

	if searchReg != "" {
		RE_SEARCH = regexp.MustCompile(searchReg)
	} else {
		RE_SEARCH = nil
	}

	var child *child_process.ChildProcess
	var err error
	if replyCommandLoop != "" {
		child, err = child_process.StartChildProcess(replyCommandLoop)
		fail("Fail to exec child", err)
	} else {
		child = nil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "irc.twitch.tv:6667")
	fail("Could not resolve TCP address", err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	fail("Could not dial TCP", err)

	_, err = conn.Write([]byte(fmt.Sprintf("PASS %v\r\nNICK %v\r\n", password, username)))
	fail("Could not write to socket for authentication", err)

	data := make([]byte, 1024)

	numRead, err := conn.Read(data)
	fail("Could not read from socket during authentication", err)

	reply := string(data[:numRead])
	if strings.Index(reply, "Login unsuccessful") > -1 {
		fail("Reply from twitch.tv", fmt.Errorf("Login failed"))
	}

	_, err = conn.Write([]byte(fmt.Sprintf("JOIN #%v\r\n", channel)))
	fail("Could not write to socket when joining channel", err)

	for {
		numRead, err := conn.Read(data)
		fail("Could not read from socket", err)

		reply := string(data[:numRead])
		for _, line := range strings.Split(reply, "\r\n") {
			c := StrToCommand(line)
			if c != nil && c.Type == TYPE_CHAT {
				if c.Username == "jtv" && ignoreAdmin == true {
				} else {
					if RE_SEARCH == nil || RE_SEARCH.MatchString(c.ChatMessage) {
						if child != nil {
							query := fmt.Sprintf("%s %s %s", c.ChatChannel, c.Username, c.ChatMessage)
							out, err := child.Query(query)
							fail("Child query failure", err)
							outString := strings.TrimSpace(out) + "\r\n"
							conn.Write([]byte(outString))
						} else {
							if replyCommand != "" {
								out, err := exec.Command(replyCommand, c.ChatChannel, c.Username, c.ChatMessage).Output()
								fail("Child query failure", err)

								outString := strings.TrimSpace(string(out)) + "\r\n"
								conn.Write([]byte(outString))
							}
						}

						fmt.Printf("%v : %s\n", c.Username, c.ChatMessage)
					}
				}
			}
		}
	}
}
