package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

type Config struct {
	Command []string
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
	err = runCommand(config.Command)
	if err != nil {
		return err
	}
	return nil
}

func runCommand(args []string) error {
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
			fmt.Print("chub | ", line)
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
