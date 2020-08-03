package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	err := run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

func run(args []string) error {
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
