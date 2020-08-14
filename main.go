package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
)

func main() {
	// Process args
	configPath := ".chub.json"
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || len(os.Args) > 2 {
			showHelp()
			os.Exit(2)
		}
		configPath = os.Args[1]
	}

	err := run(configPath)
	if err != nil {
		fmt.Println(err)
	}
}

func showHelp() {
	fmt.Println(`Usage: chub [config_path]

  Runs commands specified in the ` + "`config_path` " + `simultaneously.
  ` + "`config_path` " + `defaults to ` + "`.chub.json`" + `

  Example config:
  {
    "commands": {
      "foo": ["echo", "this is foo"],
      "bar": ["echo", "this is bar"]
    }
  }`)
}

func run(configPath string) error {
	config, err := getConfig(configPath)
	if err != nil {
		return err
	}
	err = config.Validate()
	if err != nil {
		return err
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
	commands := make([]*exec.Cmd, len(cmdNames))
	for i, name := range cmdNames {
		args := config.Commands[name]
		commands[i] = runCommand(args)
	}

	// Track commands output
	var wg sync.WaitGroup
	for i, name := range cmdNames {
		name := name
		color := getColor(i)
		prettyName := rightPad(name, maxWidth)
		prettyName += " \u258f"
		prettyName = withColor(prettyName, color)
		cmd := commands[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = trackCommand(prettyName, cmd)
			if err != nil {
				fmt.Println(prettyName + withColor(err.Error(), color))
			}
		}()
	}

	// Wait for interrupt
	shutdownSigCh := make(chan os.Signal, 1)
	signal.Notify(shutdownSigCh, syscall.SIGINT, syscall.SIGTERM)
	shutdownSig := <-shutdownSigCh
	for _, cmd := range commands {
		err := cmd.Process.Signal(shutdownSig)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Wait for commands to finish
	wg.Wait()
	return nil
}

func rightPad(s string, width int) string {
	return s + strings.Repeat(" ", width-len(s))
}

func runCommand(args []string) *exec.Cmd {
	command := args[0]
	params := args[1:]
	return exec.CommandContext(context.TODO(), command, params...)
}

func trackCommand(cmdName string, cmd *exec.Cmd) error {
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
