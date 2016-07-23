package child_process

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

type ChildProcess struct {
	Command *exec.Cmd
	In      io.WriteCloser
	Out     *bufio.Reader
}

func StartChildProcess(command string, arg ...string) (*ChildProcess, error) {
	child := &ChildProcess{}
	child.Command = exec.Command(command, arg...)

	stdout, err := child.Command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := child.Command.StdinPipe()
	if err != nil {
		return nil, err
	}

	child.In = stdin
	child.Out = bufio.NewReader(stdout)

	child.Command.Start()
	return child, nil
}

func (child *ChildProcess) Query(command string) (string, error) {
	_, err := child.In.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	out, err := child.Out.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
