package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/chzyer/readline"
)

const (
	GREEN  = "\033[32m"
	BLUE   = "\033[36m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
)

func ExecuteCommand(input string) error {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	command := args[0]

	switch command {
	case "cd":
		if len(args) < 2 {
			return errors.New("cd required a path")
		}

		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(command, args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func main() {
	//can change but I want to keep this simple
	vim := flag.Bool("vim", false, "activates vim mode")
	flag.Parse()
	username, err := user.Current()
	if err != nil {
		panic(err)
	}

	//probably not changing
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	message := strings.Builder{}

	for {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		message.Reset()
		message.WriteString(GREEN)
		message.WriteString(username.Username)
		message.WriteString("@")
		message.WriteString(hostname)
		message.WriteString(RESET)
		message.WriteString(":")
		message.WriteString(BLUE)
		message.WriteString(wd)
		message.WriteString(RESET)
		message.WriteString(YELLOW)
		message.WriteString("> ")
		message.WriteString(RESET)

		rl, err := readline.New(message.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		rl.SetVimMode(*vim)

		input, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = ExecuteCommand(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
