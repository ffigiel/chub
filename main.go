package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
)

type Config struct {
	Commands map[string][]string
}

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

	// Prepare uniform command name width
	cmdNameWidth := 0
	for cmdName := range config.Commands {
		if len(cmdName) > cmdNameWidth {
			cmdNameWidth = len(cmdName)
		}
	}

	var wg sync.WaitGroup
	for cmdName, args := range config.Commands {
		cmdName := rightPad(cmdName, cmdNameWidth)
		args := args
		wg.Add(1)
		go func() {
			err = runCommand(cmdName, args)
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
			fmt.Print(cmdName, " | ", line)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func getConfig() (Config, error) {
	var c Config
	configBytes, err := ioutil.ReadFile(".chub")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(configBytes, &c)
	return c, err
}

func readLines(ch chan<- string, r io.Reader) {
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			return
		} else if err != nil {
			fmt.Println(err)
		}
		ch <- line
	}
}
