package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}

func run() error {
	config, err := getConfig()
	if err != nil {
		return err
	}
	if len(config.Commands) == 0 {
		return fmt.Errorf("no commands found")
	}

	// Collect command names
	// also find max name length
	cmdNames := make([]string, 0, len(config.Commands))
	maxWidth := 0
	for name := range config.Commands {
		cmdNames = append(cmdNames, name)
		if len(name) > maxWidth {
			maxWidth = len(name)
		}
	}

	// Sort command names to get same colors every time
	sort.Strings(cmdNames)

	// Start the commands
	var wg sync.WaitGroup
	for i, name := range cmdNames {
		color := getColor(i)
		prettyName := rightPad(name, maxWidth)
		prettyName += " \u258f"
		prettyName = withColor(prettyName, color)
		args := config.Commands[name]
		wg.Add(1)
		go func() {
			err = runCommand(prettyName, args)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func rightPad(s string, width int) string {
	return s + strings.Repeat(" ", width-len(s))
}

func runCommand(cmdName string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}
	command := args[0]
	params := args[1:]
	cmd := exec.CommandContext(context.TODO(), command, params...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	linesCh := make(chan string, 64)
	go readLines(linesCh, stdoutReader)
	go readLines(linesCh, stderrReader)
	cmd.Start()
	go func() {
		for {
			line := <-linesCh
			fmt.Print(cmdName, line)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func readLines(ch chan<- string, r io.Reader) {
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			return
		} else if err != nil {
			fmt.Println(err)
			return
		}
		ch <- line
	}
}
