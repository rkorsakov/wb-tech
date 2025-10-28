package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Shell struct {
	reader     *bufio.Reader
	running    bool
	currentCmd *exec.Cmd
}

func NewShell() *Shell {
	return &Shell{
		reader:  bufio.NewReader(os.Stdin),
		running: true,
	}
}

func (shell *Shell) Start() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			sig := <-signalChan
			switch sig {
			case os.Interrupt:
				if shell.currentCmd != nil && shell.currentCmd.Process != nil {

					fmt.Println("\nInterrupting current command...")
					shell.currentCmd.Process.Signal(os.Interrupt)
				} else {
					fmt.Println("\nType 'exit' to quit the shell")
				}
			}
		}
	}()

	for shell.running {
		shell.printPrompt()
		input, err := shell.readInput()
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExiting...")
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		shell.executeInput(input)
	}
}

func (shell *Shell) printPrompt() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Print("$ ")
		return
	}
	fmt.Printf("%s$ ", wd)
}

func (shell *Shell) readInput() (string, error) {
	return shell.reader.ReadString('\n')
}

func (shell *Shell) executeInput(input string) {
	input = strings.TrimSpace(input)
	commands := strings.Split(input, "|")

	if len(commands) == 1 {
		shell.executeSingleCommand(commands[0])
	} else {
		shell.executePipeline(commands)
	}
}

func (shell *Shell) executeSingleCommand(cmdStr string) {
	cmdStr = strings.TrimSpace(cmdStr)
	parts := strings.Fields(cmdStr)

	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	if isBuiltinCommand(command) {
		shell.executeBuiltinCommand(command, args)
	} else {
		shell.executeExternalCommand(command, args)
	}
}

func isBuiltinCommand(cmd string) bool {
	builtins := []string{"cd", "pwd", "echo", "kill", "ps", "exit"}
	for _, builtin := range builtins {
		if cmd == builtin {
			return true
		}
	}
	return false
}

func (shell *Shell) executeBuiltinCommand(command string, args []string) {
	switch command {
	case "cd":
		shell.handleCdCommand(args)
	case "pwd":
		shell.handlePwdCommand()
	case "echo":
		shell.handleEchoCommand(args)
	case "kill":
		shell.handleKillCommand(args)
	case "ps":
		shell.handlePsCommand()
	case "exit":
		shell.running = false
	}
}

func (shell *Shell) executeExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	shell.currentCmd = cmd

	err := cmd.Run()

	shell.currentCmd = nil

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
	}
}

func (shell *Shell) executePipeline(commands []string) {
	var cmds []*exec.Cmd
	var err error

	for _, cmdStr := range commands {
		cmdStr = strings.TrimSpace(cmdStr)
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		var cmd *exec.Cmd
		if len(parts) == 1 {
			cmd = exec.Command(parts[0])
		} else {
			cmd = exec.Command(parts[0], parts[1:]...)
		}

		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		shell.currentCmd = cmds[0]
	}

	for i := 0; i < len(cmds)-1; i++ {
		if cmds[i+1].Stdin, err = cmds[i].StdoutPipe(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating pipe: %v\n", err)
			shell.currentCmd = nil
			return
		}
		cmds[i].Stderr = os.Stderr
	}

	if len(cmds) > 0 {
		cmds[len(cmds)-1].Stdout = os.Stdout
		cmds[len(cmds)-1].Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting command: %v\n", err)
			shell.currentCmd = nil
			return
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "Error waiting for command: %v\n", err)
		}
	}

	shell.currentCmd = nil
}
